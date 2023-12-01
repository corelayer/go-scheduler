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

func NewScheduler(ctx context.Context, config SchedulerConfig, catalog Catalog) *Scheduler {
	s := &Scheduler{
		config:       config,
		catalog:      catalog,
		runnerQueue:  newQueue(ctx, config.MaxSchedulableJobs),
		resultsQueue: newQueue(ctx, config.MaxSchedulableJobs),
		chRunner:     make(chan Job, config.MaxSchedulableJobs),
		chResults:    make(chan Job, config.MaxSchedulableJobs),
	}
	go s.schedule(ctx)
	go s.verifySchedule(ctx)
	go s.notifyRunners(ctx)
	return s
}

type Scheduler struct {
	config       SchedulerConfig
	catalog      Catalog
	runnerQueue  *queue
	resultsQueue *queue
	chRunner     chan Job
	chResults    chan Job
}

func (s *Scheduler) schedule(ctx context.Context) {
	time.Sleep(s.config.GetStartDelay() * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			jobs := s.catalog.Schedulable(s.config.MaxSchedulableJobs)
			if jobs != nil {
				for _, job := range jobs {
					job.Status = StatusScheduled
					s.catalog.Update(job)
					s.runnerQueue.Push(job)
				}
				// time.Sleep(s.config.GetScheduleDelay() * time.Second)
			} else {
				time.Sleep(s.config.GetNoJobsSchedulableDelay() * time.Second)
			}
		}
	}
}

func (s *Scheduler) verifySchedule(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			jobs := s.catalog.All()
			if jobs != nil {
				for _, job := range jobs {
					if job.Activate() {
						s.catalog.Activate(job.Uuid)
					}
				}
			} else {
				time.Sleep(s.config.GetScheduleDelay() * time.Second)
			}
		}
	}
}

func (s *Scheduler) notifyRunners(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job, err := s.runnerQueue.Pop()
			if err != nil {
				time.Sleep(s.config.GetNoJobsSchedulableDelay())
			}
			s.chRunner <- job
		}
	}
}

// func (s *Scheduler) Push(job *Job) {
// 	s.runnerQueue.Push(job)
// }

// func (s *Scheduler) run(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			job, err := s.runnerQueue.Pop()
// 			if err != nil {
// 				time.Sleep(1 * time.Second)
// 				continue
// 			}
// 			s.chIn <- job
// 		}
// 	}
// }

// func (s *Scheduler) process(ctx context.Context) {
// 	for {
// 		select {
// 		case job, ok := <-s.chIn:
// 			if !ok {
// 				return
// 			}
// 			for _, t := range job.Tasks {
// 				t.Execute()
// 			}
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
