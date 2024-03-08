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

	"github.com/corelayer/go-scheduler/pkg/task"
)

type OrchestratorStats struct {
	ActiveJobs  int
	EnabledJobs int
	TotalJobs   int
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
	defer o.mux.Unlock()

	enabledJobs := 0
	jobs := o.catalog.All()
	for _, job := range jobs {
		if job.IsEnabled() {
			enabledJobs++
		}
	}
	return OrchestratorStats{
		TotalJobs:   len(jobs),
		ActiveJobs:  o.activeJobs,
		EnabledJobs: enabledJobs,
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
