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

package schedule

import (
	"context"
	"strconv"
	"time"
)

func NewSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		StartDelaySeconds:             0,
		ScheduleDelaySeconds:          1,
		NoSchedulableJobsDelaySeconds: 5,
		MaxSchedulableJobs:            10,
	}
}

type SchedulerConfig struct {
	StartDelaySeconds             int
	ScheduleDelaySeconds          int
	NoSchedulableJobsDelaySeconds int
	MaxSchedulableJobs            int
}

func (c *SchedulerConfig) GetStartDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.StartDelaySeconds))
	return d
}

func (c *SchedulerConfig) GetScheduleDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.ScheduleDelaySeconds))
	return d
}

func (c *SchedulerConfig) GetNoJobsSchedulableDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.NoSchedulableJobsDelaySeconds))
	return d
}

func NewScheduler(ctx context.Context, c SchedulerConfig, r Repository) *Scheduler {
	s := &Scheduler{
		config: c,
		jobs:   r,
		queue:  NewJobQueue(ctx),
	}
	go s.schedule(ctx)
	go s.verifySchedule(ctx)
	return s
}

type Scheduler struct {
	config SchedulerConfig
	jobs   Repository
	queue  *JobQueue
}

func (s *Scheduler) schedule(ctx context.Context) {
	time.Sleep(s.config.GetStartDelay() * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			jobs := s.jobs.Schedulable(s.config.MaxSchedulableJobs)
			if jobs != nil {
				for _, job := range jobs {
					job.Status = JobStatusScheduled
					s.jobs.Update(job)
					// s.queue.Add(job)
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
			jobs := s.jobs.All()
			if jobs != nil {
				for _, job := range jobs {
					if job.Activate() {
						s.jobs.Activate(job.Uuid)
					}
				}
			} else {
				time.Sleep(s.config.GetScheduleDelay() * time.Second)
			}
		}
	}
}

// func (s *Scheduler) Add(job *Job) {
// 	s.queue.Add(job)
// }

// func (s *Scheduler) run(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		default:
// 			job, err := s.queue.Get()
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
