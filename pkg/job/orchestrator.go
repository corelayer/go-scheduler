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
)

func NewOrchestrator(catalog Catalog, config OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		config:        config,
		catalog:       catalog,
		scheduler:     NewScheduler(),
		runner:        NewRunner(),
		chSchedulable: make(chan []Job),
		chRunnable:    make(chan []Job),
		chErrors:      make(chan error),
	}
}

type Orchestrator struct {
	config    OrchestratorConfig
	catalog   Catalog
	scheduler *Scheduler
	runner    *Runner

	chSchedulable chan []Job
	chRunnable    chan []Job
	chErrors      chan error
}

func (o *Orchestrator) Run(ctx context.Context) {
	go o.handleErrors(ctx)
	go o.scheduleJobs(ctx)
	go o.runJobs(ctx)

	go o.runner.Run(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case jobs := <-o.chSchedulable:
			for _, job := range jobs {
				if err := o.scheduler.Add(job); err != nil {
					o.chErrors <- err
				}
			}
		case jobs := <-o.chRunnable:
			for _, job := range jobs {
				if err := o.runner.Add(job); err != nil {
					o.chErrors <- err
				}
			}
		}
	}
}

func (o *Orchestrator) scheduleJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			o.chSchedulable <- o.catalog.Schedulable()
			d, err := time.ParseDuration(fmt.Sprintf("%ds", o.config.ScheduleInterval))
			if err != nil {
				o.chErrors <- err
			}
			time.Sleep(d)
		}
	}
}

func (o *Orchestrator) runJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			o.chRunnable <- o.scheduler.Runnable()
		}
	}
}

func (o *Orchestrator) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-o.chErrors:
			o.config.ErrorHandler(err)
		}
	}
}

// func NewOrchestrator(ctx context.Context, c OrchestratorConfig) (*Orchestrator, error) {
// 	var (
// 		err       error
// 		scheduler *Scheduler
// 		runner    *Runner
// 	)
//
// 	chRunner := make(chan Job, c.MaxJobs)
// 	chUpdate := make(chan Job, c.MaxJobs)
//
// 	sc := NewSchedulerConfig(c.MaxJobs, c.StartupDelayMilliseconds, chRunner, chUpdate)
// 	rc := NewRunnerConfig(c.MaxJobs, c.HandlerRepository, chRunner, chUpdate)
//
// 	scheduler, err = NewScheduler(ctx, sc, c.Catalog)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	runner, err = NewRunner(ctx, rc)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	o := &Orchestrator{
// 		s:        scheduler,
// 		r:        runner,
// 		chRunner: chRunner,
// 		chUpdate: chUpdate,
// 	}
//
// 	return o, nil
// }
//
// type Orchestrator struct {
// 	s *Scheduler
// 	r *Runner
//
// 	chRunner chan Job
// 	chUpdate chan Job
// }
