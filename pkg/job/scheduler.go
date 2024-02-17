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

	"github.com/corelayer/go-scheduler/pkg/status"
)

func NewScheduler(ctx context.Context, config SchedulerConfig, catalog CatalogReadWriter) (*Scheduler, error) {
	if catalog == nil {
		return nil, fmt.Errorf("invalid catalog")
	}

	s := &Scheduler{
		config:  config,
		catalog: catalog,
	}
	go s.run(ctx)
	return s, nil
}

type Scheduler struct {
	config  SchedulerConfig
	catalog CatalogReadWriter
}

func (s *Scheduler) run(ctx context.Context) {
	queued := 0
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-s.config.chUpdate:
			if job.Status == status.StatusCompleted || job.Status == status.StatusError {
				queued--
			}
			s.catalog.UpdateActiveJob(job)
		default:
			time.Sleep(s.config.GetScheduleDelayDuration())
			jobs := s.catalog.GetActiveJobs()
			for _, job := range jobs {
				if queued < s.config.MaxJobs {
					if job.Status == status.StatusNone && job.IsDue() {
						job.SetStatus(status.StatusPending)
						s.catalog.UpdateActiveJob(job)

						s.config.chRunner <- job
						queued++
					}
				} else {
					break
				}
			}
		}
	}
}
