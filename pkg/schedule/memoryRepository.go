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

	"github.com/google/uuid"
)

func NewMemoryRepository(ctx context.Context) *MemoryRepository {
	r := MemoryRepository{
		jobs:     make([]Job, 0),
		chInput:  make(chan Job, 10),
		chDelete: make(chan uuid.UUID, 10),
		chUpdate: make(chan Job),
		mux:      sync.Mutex{},
	}
	go r.handleJobs(ctx)
	return &r
}

type MemoryRepository struct {
	jobs     []Job
	chInput  chan Job
	chDelete chan uuid.UUID
	chUpdate chan Job
	mux      sync.Mutex
}

func (r *MemoryRepository) Add(job Job) {
	r.chInput <- job
}

func (r *MemoryRepository) All() []Job {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.jobs
}

func (r *MemoryRepository) Delete(uuid uuid.UUID) {
	r.chDelete <- uuid
}

func (r *MemoryRepository) Schedulable(limit int) []Job {
	var output []Job

	r.mux.Lock()
	defer r.mux.Unlock()

	for _, job := range r.jobs {
		if job.IsSchedulable() {
			output = append(output, job)
		}
		if limit > 0 && len(output) == limit {
			break
		}
	}

	return output
}

func (r *MemoryRepository) Update(job Job) {
	r.chUpdate <- job
}

func (r *MemoryRepository) addJob(job Job) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.jobs = append(r.jobs, job)
}

func (r *MemoryRepository) deleteJob(uuid uuid.UUID) {
	r.mux.Lock()
	defer r.mux.Unlock()

	var id int

	for i, j := range r.jobs {
		if j.Uuid == uuid {
			id = i
			break
		}
	}

	r.jobs[id] = r.jobs[len(r.jobs)-1]
	r.jobs = r.jobs[:len(r.jobs)-1]
}

func (r *MemoryRepository) updateJob(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	for i, k := range r.jobs {
		if k.Uuid == job.Uuid {
			r.jobs[i] = job
			return nil
		}
	}
	return fmt.Errorf("could not update job, not found")
}

func (r *MemoryRepository) handleJobs(ctx context.Context) {
	for {
		select {
		case job, ok := <-r.chInput:
			if !ok {
				return
			}
			r.addJob(job)
		case jobId, ok := <-r.chDelete:
			if !ok {
				return
			}
			r.deleteJob(jobId)
		case job, ok := <-r.chUpdate:
			if !ok {
				return
			}
			err := r.updateJob(job)
			if err != nil {
				r.chInput <- job
			}
		case <-ctx.Done():
			return
		}
	}
}
