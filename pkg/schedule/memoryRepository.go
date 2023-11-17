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
	"sync"

	"github.com/google/uuid"
)

func NewMemoryRepository(ctx context.Context) *MemoryRepository {
	r := MemoryRepository{
		jobs:       make(map[uuid.UUID]Job, 0),
		chInput:    make(chan Job),
		chDelete:   make(chan uuid.UUID),
		chUpdate:   make(chan Job),
		chActivate: make(chan uuid.UUID),
		mux:        sync.Mutex{},
	}
	go r.handleOperations(ctx)
	return &r
}

type MemoryRepository struct {
	jobs       map[uuid.UUID]Job
	chInput    chan Job
	chDelete   chan uuid.UUID
	chUpdate   chan Job
	chActivate chan uuid.UUID
	mux        sync.Mutex
}

func (r *MemoryRepository) Activate(uuid uuid.UUID) {
	r.chActivate <- uuid
}

func (r *MemoryRepository) Add(job Job) {
	r.chInput <- job
}

func (r *MemoryRepository) All() []Job {
	r.mux.Lock()
	defer r.mux.Unlock()
	var jobs = make([]Job, 0)
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (r *MemoryRepository) Delete(uuid uuid.UUID) {
	r.chDelete <- uuid
}

// func (r *MemoryRepository) Exists(uuid uuid.UUID) bool {
// 	r.mux.Lock()
// 	defer r.mux.Unlock()
//
// 	if _, found := r.jobs[uuid]; found {
// 		return true
// 	}
// 	return false
// }

func (r *MemoryRepository) Schedulable(limit int) []Job {
	var output = make([]Job, 0)

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

func (r *MemoryRepository) activateJob(uuid uuid.UUID) {
	r.mux.Lock()
	defer r.mux.Unlock()

	job := r.jobs[uuid]
	job.Status = JobStatusSchedulable
	r.chUpdate <- job
}

func (r *MemoryRepository) addJob(job Job) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.jobs[job.Uuid] = job
}

func (r *MemoryRepository) deleteJob(uuid uuid.UUID) {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.jobs, uuid)
}

func (r *MemoryRepository) updateJob(job Job) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.jobs[job.Uuid] = job
}

func (r *MemoryRepository) handleOperations(ctx context.Context) {
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
			r.updateJob(job)
		case jobId, ok := <-r.chActivate:
			if !ok {
				return
			}
			r.activateJob(jobId)
		case <-ctx.Done():
			return
		}
	}
}
