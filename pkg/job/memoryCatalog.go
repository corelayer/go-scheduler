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
		registered: make(map[uuid.UUID]Job),
		active:     make(map[uuid.UUID]Job),
		archive:    make([]Job, 0),
		mux:        sync.Mutex{},
	}
	return &r
}

type MemoryCatalog struct {
	registered map[uuid.UUID]Job
	active     map[uuid.UUID]Job
	archive    []Job
	mux        sync.Mutex
}

// REGISTERED JOB FUNCTIONS
func (c *MemoryCatalog) CountRegisteredJobs() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.registered)
}

func (c *MemoryCatalog) CountActiveJobs() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.active)
}

func (c *MemoryCatalog) CountArchivedJobs() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(c.archive)
}

func (c *MemoryCatalog) Register(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.registered[job.Uuid] = job
	c.active[job.Uuid] = job
}

func (c *MemoryCatalog) Unregister(uuid uuid.UUID) {
	c.mux.Lock()
	defer c.mux.Unlock()

	delete(c.registered, uuid)
}

func (c *MemoryCatalog) GetArchivedJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.archive
}

// ACTIVE JOB FUNCTIONS
func (c *MemoryCatalog) GetActiveJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var jobs = make([]Job, 0)
	for _, job := range c.active {
		jobs = append(jobs, job)
	}
	return jobs
}

func (c *MemoryCatalog) GetRegisteredJobs() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var jobs = make([]Job, 0)
	for _, job := range c.registered {
		jobs = append(jobs, job)
	}
	return jobs
}

func (c *MemoryCatalog) UpdateActiveJob(job Job) {
	c.mux.Lock()
	defer c.mux.Unlock()

	// fmt.Printf("Updating job \"%s\" - status \"%s\"\r\n", job.Name, job.Status)
	if job.Status == StatusCompleted || job.Status == StatusError {
		// Copy job to archive
		c.archive = append(c.archive, job)

		// Delete job from active jobs
		delete(c.active, job.Uuid)

		// Reactivate job if the job must be repeated, delete from repository if not
		if c.registered[job.Uuid].Repeat {
			// j := c.registered[job.Uuid]
			// task := PrintTask{Message: "##### " + strconv.Itoa(time.Now().Minute())}
			// j.Tasks.Tasks[1] = task
			// fmt.Printf("### Adding repeating job %s - status %s\r\n", j.Name, j.Status)
			// c.active[j.Uuid] = j
			c.active[job.Uuid] = c.registered[job.Uuid]
		} else {
			delete(c.registered, job.Uuid)
		}
	} else {
		c.active[job.Uuid] = job
	}
}
