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

	"github.com/google/uuid"
)

func NewMemoryCatalog() *MemoryCatalog {
	r := MemoryCatalog{
		jobs: make(map[uuid.UUID]Job),
		mux:  sync.Mutex{},
	}
	return &r
}

type MemoryCatalog struct {
	jobs map[uuid.UUID]Job
	mux  sync.Mutex
}

func (c *MemoryCatalog) GetNotSchedulableJobs() []Job {
	var jobs = make([]Job, 0)

	c.mux.Lock()
	defer c.mux.Unlock()
	for _, job := range c.jobs {
		if !job.IsSchedulable() {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

func (c *MemoryCatalog) Count() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return len(c.jobs)
}

func (c *MemoryCatalog) GetDueJobs(limit int) []Job {
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

func (c *MemoryCatalog) GetRunnableJobs(limit int) []Job {
	var output = make([]Job, 0)

	c.mux.Lock()
	defer c.mux.Unlock()

	for _, job := range c.jobs {
		if job.IsRunnable() {
			output = append(output, job)
		}
		if limit > 0 && len(output) == limit {
			break
		}
	}

	return output
}

func (c *MemoryCatalog) Activate(uuid uuid.UUID) {
	c.mux.Lock()
	defer c.mux.Unlock()
	job := c.jobs[uuid]
	job.Status = StatusIsDue
	c.jobs[uuid] = job
}

func (c *MemoryCatalog) Add(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.jobs[job.Uuid] = job
}

func (c *MemoryCatalog) Delete(uuid uuid.UUID) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.jobs, uuid)
}

func (c *MemoryCatalog) Update(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.jobs[job.Uuid] = job
}
