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
	"fmt"
	"sync"
)

func NewJobQueue(ctx context.Context) *JobQueue {
	s := &JobQueue{
		jobs: make([]Job, 0),
		chIn: make(chan Job, 10),
	}
	go s.handleJobInput(ctx)
	return s
}

type JobQueue struct {
	jobs []Job
	chIn chan Job

	mux sync.Mutex
}

func (s *JobQueue) Length() int {
	return len(s.jobs)
}
func (s *JobQueue) Capacity() int {
	return cap(s.jobs)
}

func (s *JobQueue) Add(job Job) {
	s.chIn <- job
}

func (s *JobQueue) Get() (Job, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if len(s.jobs) > 0 {
		job := s.jobs[0]
		s.jobs = s.jobs[1:]
		return job, nil
	}
	return Job{}, fmt.Errorf("no jobs available")
}

func (s *JobQueue) handleJobInput(ctx context.Context) {
	for {
		select {
		case job, ok := <-s.chIn:
			if !ok {
				return
			}
			s.mux.Lock()
			s.jobs = append(s.jobs, job)
			s.mux.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
