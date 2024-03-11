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
	"sync"
	"time"

	"github.com/corelayer/go-scheduler/pkg/task"
)

func NewOrchestrator(catalog Catalog, taskHandlers *task.HandlerRepository, config OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		config:       config,
		catalog:      catalog,
		taskHandlers: taskHandlers,
		chRunnerIn:   make(chan Job, config.MaxJobs),
		chMessages:   make(chan task.IntercomMessage),
		chErrors:     make(chan error),
		runningJobs:  0,
		mux:          sync.Mutex{},
	}
}

type Orchestrator struct {
	config       OrchestratorConfig
	catalog      Catalog
	taskHandlers *task.HandlerRepository
	chRunnerIn   chan Job
	chMessages   chan task.IntercomMessage
	chErrors     chan error
	runningJobs  int
	isStarted    bool
	mux          sync.Mutex
}

func (o *Orchestrator) IsStarted() bool {
	o.mux.Lock()
	defer o.mux.Unlock()

	return o.isStarted
}

func (o *Orchestrator) Start(ctx context.Context) {
	o.mux.Lock()
	time.Sleep(o.config.StartDelay)

	// Make sure the orchestrator is ready to handle errors and messages before launching any other goroutine
	go o.handleErrors()
	go o.handleMessages()

	// Launch goroutines in the order of the "normal" job flow
	go o.handleInactiveJobs(ctx)
	go o.handleAvailableJobs(ctx)
	go o.handleSchedulableJobs(ctx)
	go o.handleRunnableJobs(ctx)
	go o.handlePendingJobs(ctx)

	for i := 0; i < o.config.MaxJobs; i++ {
		go o.handleActiveJob()
	}

	go o.handleShutdown(ctx)

	o.isStarted = true
	o.mux.Unlock()
}

func (o *Orchestrator) Statistics() OrchestratorStats {
	o.mux.Lock()
	jobs := o.catalog.All()
	runningJobs := o.runningJobs
	o.mux.Unlock()

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
	taskStats := make([]TaskStats, 0)

	totalTasks := 0
	completedTasks := 0

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

		totalTasks += job.Tasks.Count()

		currentResult := job.CurrentResult()
		completedTasks += len(currentResult.Tasks)

		hasErrors := false
		for _, t := range currentResult.Tasks {
			if t.Status() == task.StatusError || t.Status() == task.StatusCanceled {
				hasErrors = true
				break
			}
		}
		taskStats = append(taskStats, TaskStats{Uuid: job.Uuid, Name: job.Name, Completed: float64(len(currentResult.Tasks)), Total: float64(job.Tasks.Count()), HasErrors: hasErrors})
	}
	return OrchestratorStats{
		Job: GlobalStats{
			ConfiguredJobs:  float64(configuredJobs),
			EnabledJobs:     float64(enabledJobs),
			DisabledJobs:    float64(disabledJobs),
			InactiveJobs:    float64(inactiveJobs),
			AvailableJobs:   float64(availableJobs),
			SchedulableJobs: float64(schedulableJobs),
			RunnableJobs:    float64(runnableJobs),
			PendingJobs:     float64(pendingJobs),
			ActiveJobs:      float64(activeJobs),
			RunningJobs:     float64(runningJobs),
			CompletedTasks:  float64(completedTasks),
			TotalTasks:      float64(totalTasks),
		},
		Tasks: taskStats,
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

func (o *Orchestrator) handleActiveJob() {
	for {
		job, ok := <-o.chRunnerIn
		if !ok {
			return
		}
		o.runningJobsIncrease()

		// Update job data
		result := Result{
			Start:  time.Now(),
			Status: StatusActive,
		}
		job.AddResult(result)

		// Send job update to catalog, so we can track active jobs
		if err := o.catalog.Update(job); err != nil {
			o.chErrors <- err
		}

		job.Tasks.active = true

		// Run all task for job
		intercom := task.NewIntercom(job.Name, o.chMessages)
		pipeline := make(chan *task.Pipeline, 1)
		pipeline <- &task.Pipeline{Intercom: intercom, Data: make(map[string]interface{})}

		for i, t := range job.Tasks.All() {
			job.Tasks.activeIdx = i
			job.Tasks.executed = append(job.Tasks.executed, t)

			// Execute current task
			taskResult := o.taskHandlers.Execute(t, pipeline)

			job.Tasks.executed[job.Tasks.activeIdx] = taskResult

			result.Tasks = job.Tasks.Executed()
			job.UpdateResult(result)

			if err := o.catalog.Update(job); err != nil {
				o.chErrors <- err
			}
		}
		close(pipeline)

		job.Tasks.active = false

		result.Finish = time.Now()
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

		o.runningJobsDecrease()
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

func (o *Orchestrator) handlePendingJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(o.chRunnerIn)
			return
		default:
			for _, job := range o.catalog.PendingJobs() {
				job.SetStatus(StatusActive)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
				} else {
					o.chRunnerIn <- job
				}
			}
		}
	}
}

func (o *Orchestrator) handleRunnableJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, job := range o.catalog.RunnableJobs() {
				job.SetStatus(StatusPending)
				if err := o.catalog.Update(job); err != nil {
					o.chErrors <- err
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

			time.Sleep(o.config.ScheduleInterval)
		}
	}
}

func (o *Orchestrator) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	for {
		o.mux.Lock()
		if o.runningJobs != 0 {
			o.mux.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		} else {
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

func (o *Orchestrator) runningJobsIncrease() {
	o.mux.Lock()
	o.runningJobs++
	o.mux.Unlock()
}

func (o *Orchestrator) runningJobsDecrease() {
	o.mux.Lock()
	o.runningJobs--
	o.mux.Unlock()
}
