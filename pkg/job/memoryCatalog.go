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
	"sync"

	"github.com/google/uuid"
)

func NewMemoryCatalog(ctx context.Context) *MemoryCatalog {
	r := MemoryCatalog{
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

type MemoryCatalog struct {
	jobs       map[uuid.UUID]Job
	chInput    chan Job
	chDelete   chan uuid.UUID
	chUpdate   chan Job
	chActivate chan uuid.UUID
	mux        sync.Mutex
}

func (c *MemoryCatalog) Activate(uuid uuid.UUID) {
	c.chActivate <- uuid
}

func (c *MemoryCatalog) Add(job Job) {
	c.chInput <- job
}

func (c *MemoryCatalog) All() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()
	var jobs = make([]Job, 0)
	for _, job := range c.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (c *MemoryCatalog) Delete(uuid uuid.UUID) {
	c.chDelete <- uuid
}

// func (r *MemoryCatalog) Exists(uuid uuid.UUID) bool {
// 	r.mux.Lock()
// 	defer r.mux.Unlock()
//
// 	if _, found := r.catalog[uuid]; found {
// 		return true
// 	}
// 	return false
// }

func (c *MemoryCatalog) Schedulable(limit int) []Job {
	var output = make([]Job, 0)

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, job := range c.jobs {
		if job.IsSchedulable() {
			output = append(output, job)
		}
		if limit > 0 && len(output) == limit {
			break
		}
	}

	return output
}

func (c *MemoryCatalog) Update(job Job) {
	c.chUpdate <- job
}

func (c *MemoryCatalog) activateJob(uuid uuid.UUID) {
	c.mux.Lock()
	defer c.mux.Unlock()

	job := c.jobs[uuid]
	job.Status = StatusSchedulable
	c.chUpdate <- job
}

func (c *MemoryCatalog) addJob(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.jobs[job.Uuid] = job
}

func (c *MemoryCatalog) deleteJob(uuid uuid.UUID) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.jobs, uuid)
}

func (c *MemoryCatalog) updateJob(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.jobs[job.Uuid] = job
}

func (c *MemoryCatalog) handleOperations(ctx context.Context) {
	for {
		select {
		case job, ok := <-c.chInput:
			if !ok {
				return
			}
			c.addJob(job)
		case jobId, ok := <-c.chDelete:
			if !ok {
				return
			}
			c.deleteJob(jobId)
		case job, ok := <-c.chUpdate:
			if !ok {
				return
			}
			c.updateJob(job)
		case jobId, ok := <-c.chActivate:
			if !ok {
				return
			}
			c.activateJob(jobId)
		case <-ctx.Done():
			return
		}
	}
}
