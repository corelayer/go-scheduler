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
	"sync"
)

func NewScheduler() *Scheduler {
	return &Scheduler{
		repository: NewRepository(),
		mux:        sync.Mutex{},
	}
}

type Scheduler struct {
	repository *Repository
	mux        sync.Mutex
}

func (s *Scheduler) Add(job Job) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.repository.Add(job)
}

func (s *Scheduler) Runnable() []Job {
	s.mux.Lock()
	defer s.mux.Unlock()

	var jobs = make([]Job, 0)
	available := s.repository.All()
	for _, job := range available {
		if job.IsRunnable() {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

// func NewScheduler(ctx context.Context, config SchedulerConfig, catalog CatalogReadWriter) (*Scheduler, error) {
// 	if catalog == nil {
// 		return nil, fmt.Errorf("invalid catalog")
// 	}
//
// 	s := &Scheduler{
// 		config:  config,
// 		catalog: catalog,
// 	}
// 	go s.run(ctx)
// 	return s, nil
// }
//
// type Scheduler struct {
// 	config  SchedulerConfig
// 	catalog CatalogReadWriter
// }
//
// func (s *Scheduler) run(ctx context.Context) {
// 	time.Sleep(s.config.GetStartupDelayDuration())
// 	queued := 0
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case job := <-s.config.chUpdate:
// 			if job.Status == StatusCompleted || job.Status == StatusError {
// 				queued--
// 			}
// 			s.catalog.UpdateActiveJob(job)
// 		default:
// 			time.Sleep(s.config.GetScheduleDelayDuration())
// 			jobs := s.catalog.GetActiveJobs()
// 			for _, job := range jobs {
// 				if queued < s.config.MaxJobs {
// 					if job.Status == StatusNone && job.IsDue() {
// 						job.SetStatus(StatusPending)
// 						s.catalog.UpdateActiveJob(job)
//
// 						s.config.chRunner <- job
// 						queued++
// 					}
// 				} else {
// 					break
// 				}
// 			}
// 		}
// 	}
// }
