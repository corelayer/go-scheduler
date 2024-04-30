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
	return &MemoryCatalog{
		jobs: make(map[uuid.UUID]Job, 0),
		mux:  sync.Mutex{},
	}
}

type MemoryCatalog struct {
	jobs map[uuid.UUID]Job
	mux  sync.Mutex
}

func (c *MemoryCatalog) Add(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	for k := range c.jobs {
		if job.Uuid == k {
			return ErrExist
		}
	}

	c.jobs[job.Uuid] = job
	return nil
}

func (c *MemoryCatalog) All() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	jobs := make([]Job, 0)
	for _, job := range c.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

func (c *MemoryCatalog) AvailableJobs() []Job {
	return c.GetJobsByStatus(StatusAvailable)
}

func (c *MemoryCatalog) Count() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.jobs)
}

func (c *MemoryCatalog) Delete(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[jobId]; !found {
		return ErrNotFound
	}
	delete(c.jobs, jobId)
	return nil
}

func (c *MemoryCatalog) Disable(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[jobId]; !found {
		return ErrNotFound
	}
	job := c.jobs[jobId]
	job.Disable()
	c.jobs[jobId] = job

	return nil
}

func (c *MemoryCatalog) Enable(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[jobId]; !found {
		return ErrNotFound
	}
	job := c.jobs[jobId]
	job.Enable()
	c.jobs[jobId] = job

	return nil
}

func (c *MemoryCatalog) Exists(id uuid.UUID) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[id]; found {
		return true
	}
	return false
}

func (c *MemoryCatalog) Get(id uuid.UUID) (Job, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[id]; !found {
		return Job{}, ErrNotFound
	}
	return c.jobs[id], nil
}

func (c *MemoryCatalog) GetJobsByStatus(status Status) []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var jobs = make([]Job, 0)
	for _, job := range c.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}
	return jobs

}

func (c *MemoryCatalog) HasEnabledJobs() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, job := range c.jobs {
		if job.IsEnabled() {
			return true
		}
	}
	return false
}

func (c *MemoryCatalog) InactiveJobs() []Job {
	return c.GetJobsByStatus(StatusInactive)
}

func (c *MemoryCatalog) PendingJobs() []Job {
	return c.GetJobsByStatus(StatusPending)
}

func (c *MemoryCatalog) RunnableJobs() []Job {
	return c.GetJobsByStatus(StatusRunnable)
}

func (c *MemoryCatalog) SchedulableJobs() []Job {
	return c.GetJobsByStatus(StatusSchedulable)
}

func (c *MemoryCatalog) Update(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.jobs[job.Uuid]; !found {
		return ErrNotFound
	}
	c.jobs[job.Uuid] = job
	return nil
}
