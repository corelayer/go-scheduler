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
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func NewMemoryCatalog() *MemoryCatalog {
	return &MemoryCatalog{
		repository: NewRepository(),
		mux:        sync.Mutex{},
	}
}

type MemoryCatalog struct {
	repository *Repository
	mux        sync.Mutex
}

func (c *MemoryCatalog) Add(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.repository.Add(job)
}

func (c *MemoryCatalog) All() map[uuid.UUID]Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.repository.GetAll()
}

func (c *MemoryCatalog) AvailableJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var available = make([]Job, 0)

	jobs := c.repository.GetAll()
	for _, job := range jobs {
		if job.IsAvailable() {
			available = append(available, job)
		}
	}
	return available
}

func (c *MemoryCatalog) Delete(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.repository.Delete(jobId)
}

func (c *MemoryCatalog) Disable(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	job, err := c.repository.Get(jobId)
	if err != nil {
		return err
	}

	job.Disable()
	return c.repository.Update(job)
}

func (c *MemoryCatalog) Enable(jobId uuid.UUID) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	job, err := c.repository.Get(jobId)
	if err != nil {
		return err
	}

	job.Enable()
	return c.repository.Update(job)
}

func (c *MemoryCatalog) HasEnabledJobs() bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	var enabledCount int
	var disabledCount int
	for _, job := range c.repository.GetAll() {
		if job.IsEnabled() {
			enabledCount++
		} else {
			disabledCount++
		}
	}
	fmt.Println(enabledCount, disabledCount, enabledCount != 0)
	return enabledCount != 0
}

func (c *MemoryCatalog) InactiveJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var inactive = make([]Job, 0)
	for _, job := range c.repository.GetAll() {
		if job.IsInactive() {
			inactive = append(inactive, job)
		}
	}
	return inactive
}

func (c *MemoryCatalog) RunnableJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var runnable = make([]Job, 0)
	jobs := c.repository.GetAll()
	for _, job := range jobs {
		if job.IsRunnable() {
			runnable = append(runnable, job)
		}
	}
	return runnable
}

func (c *MemoryCatalog) SchedulableJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var schedulable = make([]Job, 0)
	jobs := c.repository.GetAll()
	for _, job := range jobs {
		if job.IsSchedulable() {
			schedulable = append(schedulable, job)
		}
	}
	return schedulable
}

func (c *MemoryCatalog) Update(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if !job.IsInactive() {
		job.Disable()
	}

	return c.repository.Update(job)
}
