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

func NewRunner(ctx context.Context, config RunnerConfig, catalog CatalogReadWriter) (*Runner, error) {
	if catalog == nil {
		return nil, fmt.Errorf("invalid catalog")
	}

	r := &Runner{
		config:        config,
		catalog:       catalog,
		queue:         NewMemoryQueue(),
		chWorkerInput: make(chan Job),
	}

	workers := make([]*Worker, config.maxConcurrentJobs)
	for i := 0; i < config.maxConcurrentJobs; i++ {
		workers[i] = NewWorker(ctx, WorkerConfig{
			id:                        i,
			idleSleepTimeMilliseconds: 10,
		}, r.chWorkerInput)
	}
	r.workers = workers

	go r.queueJobs(ctx)
	go r.runJobs(ctx)
	return r, nil

}

type Runner struct {
	config        RunnerConfig
	catalog       CatalogReadWriter
	queue         Queue
	workers       []*Worker
	chWorkerInput chan Job
}

func (r *Runner) queueJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ql := r.queue.Length()
			if ql >= r.config.maxConcurrentJobs {
				time.Sleep(r.config.GetIdleDelayDuration())
				continue
			}

			jobs := r.catalog.GetRunnableJobs(1)
			if len(jobs) == 0 {
				time.Sleep(r.config.GetIdleDelayDuration())
				continue
			}
			for _, job := range jobs {
				job.Status = StatusScheduled
				r.catalog.Update(job)
				// Add job to worker queue
				r.queue.Push(job)
			}

		}
	}
}

func (r *Runner) runJobs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := r.queue.Pop()
			if err != nil {
				time.Sleep(r.config.GetIdleDelayDuration())
				continue
			}

			job.Status = StatusPending
			r.catalog.Update(job)

			sent := false
			for {
				if sent {
					break
				}
				select {
				case r.chWorkerInput <- job:
					sent = true
				}
			}
		}
	}
}
