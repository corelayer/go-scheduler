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

func NewScheduler(ctx context.Context, config SchedulerConfig, catalog CatalogReadWriter) (*Scheduler, error) {
	if catalog == nil {
		return nil, fmt.Errorf("invalid catalog")
	}

	s := &Scheduler{
		config:  config,
		catalog: catalog,
	}
	go s.schedule(ctx)
	go s.verifySchedule(ctx)
	return s, nil
}

type Scheduler struct {
	config  SchedulerConfig
	catalog CatalogReadWriter
}

func (s *Scheduler) schedule(ctx context.Context) {
	time.Sleep(s.config.GetStartDelayDuration())
	for {
		select {
		case <-ctx.Done():
			return
		default:
			jobs := s.catalog.GetDueJobs(s.config.maxSchedulableJobs)
			if len(jobs) == 0 {
				time.Sleep(s.config.GetIdleDelayDuration())
				continue
			}
			for _, job := range jobs {
				job.Status = StatusSchedulable
				s.catalog.Update(job)
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
			time.Sleep(s.config.GetScheduleDelayDuration())
			jobs := s.catalog.GetNotSchedulableJobs()
			if len(jobs) == 0 {
				time.Sleep(s.config.GetIdleDelayDuration())
				continue
			}
			for _, job := range jobs {
				if job.Status == StatusNone && job.IsDue() {
					s.catalog.Activate(job.Uuid)
				}
			}
		}
	}
}
