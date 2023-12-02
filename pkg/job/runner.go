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
	"time"
)

func NewRunner(ctx context.Context, config RunnerConfig, catalog CatalogReadWriter) *Runner {
	r := &Runner{
		config:        config,
		catalog:       catalog,
		queue:         NewMemoryQueue(),
		chWorkerInput: make(chan Job),
		// chWorkerOutput: make(chan Job, config.MaxConcurrentJobs),
	}

	workers := make([]*Worker, config.MaxConcurrentJobs)
	for i := 0; i < config.MaxConcurrentJobs; i++ {
		workers[i] = NewWorker(ctx, i, r.chWorkerInput)
	}
	r.workers = workers

	go r.queueJobs(ctx)
	go r.runJobs(ctx)
	return r

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
			if ql >= r.config.MaxConcurrentJobs {
				// fmt.Println("Waiting for space in the queue:", ql)
				time.Sleep(r.config.GetNoRunnableJobsDelayDuration())
				continue
			}

			jobs := r.catalog.GetRunnableJobs(1)
			if len(jobs) == 0 {
				// fmt.Println("Waiting for runnable jobs in catalog")
				time.Sleep(r.config.GetNoRunnableJobsDelayDuration())
				continue
			}
			for _, job := range jobs {
				job.Status = StatusScheduled
				r.catalog.Update(job)
				// fmt.Println("Adding job to queue:", job.Name)
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
				time.Sleep(r.config.GetNoRunnableJobsDelayDuration())
				continue
			}

			job.Status = StatusPending
			r.catalog.Update(job)

			sent := false
			for {
				if sent {
					// fmt.Printf("Job %s sent to worker\r\n", job.Name)
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
