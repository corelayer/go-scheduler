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
	"time"

	"github.com/corelayer/go-scheduler/pkg/task"
)

func NewOrchestrator(catalog Catalog, taskHandlers *task.HandlerRepository, config OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		config:       config,
		catalog:      catalog,
		taskHandlers: taskHandlers,
		chRunnerIn:   make(chan Job),
		chRunnerOut:  make(chan Job),
		chMessages:   make(chan task.IntercomMessage),
		chErrors:     make(chan error),
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
}

func (o *Orchestrator) Start(ctx context.Context) {
	go o.handleErrors(ctx)
	go o.handleMessages(ctx)
	go o.handleInactiveJobs(ctx)
	go o.handleAvailableJobs(ctx)
	go o.handleSchedulableJobs(ctx)
	go o.handleRunnableJobs(ctx)
	go o.handleResults(ctx)

	for i := 0; i < o.config.MaxJobs; i++ {
		go o.handleJobs(ctx)
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

func (o *Orchestrator) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-o.chErrors:
			if !ok {
				return
			}

			if o.config.ErrorHandler != nil {
				o.config.ErrorHandler(err)
			}
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

func (o *Orchestrator) handleJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-o.chRunnerIn:
			if !ok {
				return
			}

			// Update job data
			result := Result{
				start:   time.Now(),
				RunTime: 0,
				Status:  StatusActive,
			}
			job.AddResult(result)
			job.SetStatus(StatusActive)
			o.chRunnerOut <- job

			// Run all task for job
			intercom := task.NewIntercom(job.Name, o.chMessages)
			job.Tasks.Execute(o.taskHandlers, intercom)

			result.finish = time.Now()
			result.RunTime = result.finish.Sub(result.start)
			result.messages = intercom.GetAll()
			if intercom.HasErrors() {
				result.Status = StatusError
			} else {
				result.Status = StatusCompleted
			}
			job.UpdateResult(result)
			job.SetStatus(result.Status)
			o.chRunnerOut <- job
		}
	}
}

func (o *Orchestrator) handleMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-o.chMessages:
			if !ok {
				return
			}

			if o.config.MessageHandler != nil {
				o.config.MessageHandler(message)

			}
		}
	}
}

func (o *Orchestrator) handleResults(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-o.chRunnerOut:
			if !ok {
				return
			}

			if !job.IsActive() {
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
