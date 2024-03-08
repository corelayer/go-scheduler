/*
 * Copyright 2023 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/task"
)

type JobStats struct {
	ConfiguredJobs  float64
	EnabledJobs     float64
	DisabledJobs    float64
	ActiveJobs      float64
	AvailableJobs   float64
	InactiveJobs    float64
	PendingJobs     float64
	RunnableJobs    float64
	SchedulableJobs float64
}

type TaskStats struct {
	Uuid      uuid.UUID
	Name      string
	Completed float64
	Total     float64
}

type OrchestratorStats struct {
	Job   JobStats
	Tasks []TaskStats
}

func NewOrchestrator(catalog Catalog, taskHandlers *task.HandlerRepository, config OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		config:       config,
		catalog:      catalog,
		taskHandlers: taskHandlers,
		chRunnerIn:   make(chan Job, config.MaxJobs),
		chRunnerOut:  make(chan Job),
		chMessages:   make(chan task.IntercomMessage),
		chErrors:     make(chan error),
		activeJobs:   0,
		mux:          sync.Mutex{},
	}
}

type Orchestrator struct {
	config       OrchestratorConfig
	catalog      Catalog
	taskHandlers *task.HandlerRepository
	chRunnerIn   chan Job
	chRunnerOut  chan Job
	chMessages   chan task.IntercomMessage
	chErrors     chan error
	activeJobs   int
	isStarted    bool
	mux          sync.Mutex
}

func (o *Orchestrator) IsStarted() bool {
	o.mux.Lock()
	defer o.mux.Unlock()

	return o.isStarted
}

func (o *Orchestrator) Start(ctx context.Context) {
	// Make sure the orchestrator is ready to handle errors and messages before launching any other goroutine
	go o.handleErrors()
	go o.handleMessages()

	// Launch goroutines in the order of the "normal" job flow
	go o.handleInactiveJobs(ctx)
	go o.handleAvailableJobs(ctx)
	go o.handleSchedulableJobs(ctx)
	go o.handleRunnableJobs(ctx)
	go o.handleResults()

	for i := 0; i < o.config.MaxJobs; i++ {
		go o.handleJobs()
	}

	go o.handleShutdown(ctx)

	o.mux.Lock()
	o.isStarted = true
	o.mux.Unlock()
}

func (o *Orchestrator) Statistics() OrchestratorStats {
	o.mux.Lock()
	jobs := o.catalog.All()
	defer o.mux.Unlock()

	configuredJobs := len(jobs)
	enabledJobs := 0
	disabledJobs := 0
	activeJobs := 0
	availableJobs := 0
	inactiveJobs := 0
	pendingJobs := 0
	runnableJobs := 0
	schedulableJobs := 0
	finishedJobs := 0
	completedTasks := make([]TaskStats, 0)

	for _, job := range jobs {
		switch job.IsEnabled() {
		case true:
			enabledJobs++
		case false:
			disabledJobs++
		}

		switch job.Status {
		case StatusActive:
			activeJobs++
		case StatusAvailable:
			availableJobs++
		case StatusInactive:
			inactiveJobs++
			finishedJobs += job.CountRuns()
		case StatusPending:
			pendingJobs++
		case StatusRunnable:
			runnableJobs++
		case StatusSchedulable:
			schedulableJobs++
		default:
		}

		completedTasks = append(completedTasks, TaskStats{Uuid: job.Uuid, Name: job.Name, Completed: float64(len(job.CurrentResult().Tasks)), Total: float64(job.Tasks.Count())})
	}
	return OrchestratorStats{
		Job: JobStats{
			ConfiguredJobs:  float64(configuredJobs),
			EnabledJobs:     float64(enabledJobs),
			DisabledJobs:    float64(disabledJobs),
			ActiveJobs:      float64(activeJobs),
			AvailableJobs:   float64(availableJobs),
			InactiveJobs:    float64(inactiveJobs),
			PendingJobs:     float64(pendingJobs),
			RunnableJobs:    float64(runnableJobs),
			SchedulableJobs: float64(schedulableJobs),
		},
		Tasks: completedTasks,
	}
}

func (o *Orchestrator) handleAvailableJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range o.catalog.AvailableJobs() {
				job.SetStatus(StatusSchedulable)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
				}
			}
		}
	}
}

func (o *Orchestrator) handleErrors() {
	for {
		err, ok := <-o.chErrors
		if !ok {
			return
		}
		if o.config.ErrorHandler != nil {
			o.config.ErrorHandler(err)
		}
	}
}

func (o *Orchestrator) handleInactiveJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range o.catalog.InactiveJobs() {
				job.SetStatus(StatusAvailable)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
				}
			}
		}
	}
}

func (o *Orchestrator) handleJobs() {
	for {
		job, ok := <-o.chRunnerIn
		if !ok {
			return
		}

		o.mux.Lock()
		o.activeJobs++
		o.mux.Unlock()

		// Update job data
		result := Result{
			Start:   time.Now(),
			RunTime: 0,
			Status:  StatusActive,
		}
		job.AddResult(result)
		job.SetStatus(StatusActive)
		// Send job update to catalog, so we can track active jobs
		o.chRunnerOut <- job

		// Run all task for job
		intercom := task.NewIntercom(job.Name, o.chMessages)
		job.Tasks.Execute(o.taskHandlers, intercom)

		result.Finish = time.Now()
		result.RunTime = result.Finish.Sub(result.Start)
		result.Tasks = job.Tasks.Executed()
		result.Messages = intercom.GetAll()
		if intercom.HasErrors() {
			result.Status = StatusError
		} else {
			result.Status = StatusCompleted
		}
		job.UpdateResult(result)
		job.Tasks.ResetHistory()
		job.SetStatus(result.Status)

		o.chRunnerOut <- job

		o.mux.Lock()
		o.activeJobs--
		o.mux.Unlock()
	}
}

func (o *Orchestrator) handleMessages() {
	for {
		message, ok := <-o.chMessages
		if !ok {
			return
		}
		if o.config.MessageHandler != nil {
			o.config.MessageHandler(message)

		}
	}
}

func (o *Orchestrator) handleResults() {
	for {
		job, ok := <-o.chRunnerOut
		if !ok {
			return
		}

		// Job status is not active --> StatusCompleted or StatusError
		if !job.IsActive() {
			// Disable job if it does not need to be run again
			if !job.IsEligible() {
				job.Disable()
			} else {
				job.SetStatus(StatusInactive)
			}
		}

		if err := o.catalog.Update(job); err != nil {
			o.chErrors <- err
		}
	}
}

func (o *Orchestrator) handleRunnableJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(o.chRunnerIn)
			return
		default:
			for _, job := range o.catalog.RunnableJobs() {
				job.SetStatus(StatusPending)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
				} else {
					o.chRunnerIn <- job
				}
			}
		}
	}
}

func (o *Orchestrator) handleSchedulableJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range o.catalog.SchedulableJobs() {
				job.SetStatus(StatusRunnable)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
				}
			}
			d, err := time.ParseDuration(fmt.Sprintf("%dms", o.config.ScheduleInterval))
			if err != nil {
				o.chErrors <- err
			}
			time.Sleep(d)
		}
	}
}

func (o *Orchestrator) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	for {
		o.mux.Lock()
		if o.activeJobs != 0 {
			o.mux.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
			close(o.chRunnerOut)
			close(o.chMessages)
			close(o.chErrors)
			o.mux.Unlock()
			break
		}
	}
	o.mux.Lock()
	o.isStarted = false
	o.mux.Unlock()
}
