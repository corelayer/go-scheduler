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

type Catalog interface {
	Add(job Job) error
	Delete(jobId uuid.UUID) error
	AddResult(job Job) error
	Result(jobId uuid.UUID, runId uuid.UUID) (Job, error)
	Results(jobId uuid.UUID) (map[uuid.UUID]Job, error)
	RunCount(jobId uuid.UUID) (int, error)
	Schedulable() []Job
	Update(job Job) error
}

func NewMemoryCatalog() *MemoryCatalog {
	return &MemoryCatalog{
		repository: NewRepository(),
		results:    make(map[uuid.UUID]*Repository),
		mux:        sync.Mutex{},
	}
}

type MemoryCatalog struct {
	repository *Repository
	results    map[uuid.UUID]*Repository
	mux        sync.Mutex
}

func (c *MemoryCatalog) Add(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.repository.Add(job)
}

func (c *MemoryCatalog) AddResult(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.results[job.id]; !found {
		c.results[job.id] = NewRepository()
	}

	return c.results[job.id].Add(job)
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

func (c *MemoryCatalog) Result(jobId uuid.UUID, runId uuid.UUID) (Job, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.results[jobId]; !found {
		return Job{}, ErrNotFound
	}

	return c.results[jobId].Get(runId)
}

func (c *MemoryCatalog) Results(jobId uuid.UUID) (map[uuid.UUID]Job, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, found := c.results[jobId]; !found {
		return nil, ErrNotFound
	}

	return c.results[jobId].All(), nil
}

func (c *MemoryCatalog) RunCount(jobId uuid.UUID) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if !c.repository.Exists(jobId) {
		return 0, ErrNotFound
	}

	if _, found := c.results[jobId]; !found {
		return 0, nil
	}

	return c.results[jobId].Count(), nil
}

func (c *MemoryCatalog) Schedulable() []Job {
	c.mux.Lock()
	defer c.mux.Unlock()

	var jobs = make([]Job, 0)
	available := c.repository.All()
	for _, job := range available {
		if !job.Enabled() {
			continue
		}

		if job.maxRuns == 0 {
			jobs = append(jobs, job)
			continue
		}

		if c.results[job.id].Count() < job.maxRuns {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

func (c *MemoryCatalog) Update(job Job) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.Update(job)
}

// func NewMemoryCatalog() *MemoryCatalog {
// 	r := MemoryCatalog{
// 		// registered: make(map[uuid.UUID]Job),
// 		// active:     make(map[uuid.UUID]Job),
// 		// archive:    make([]Job, 0),
// 		mux: sync.Mutex{},
//
// 		repository: make(map[uuid.UUID]Definition),
// 		// schedulable: make(map[uuid.UUID]Schedulable),
// 		// activated:   make(map[uuid.UUID]Active),
// 		archived: make(map[uuid.UUID][]Result),
// 	}
// 	return &r
// }
//
// type MemoryCatalog struct {
// 	// registered map[uuid.UUID]Job
// 	// active     map[uuid.UUID]Job
// 	// archive    []Job
// 	mux sync.Mutex
//
// 	repository  map[uuid.UUID]Definition
// 	schedulable map[uuid.UUID]Schedulable
// 	activated   map[uuid.UUID]Active
// 	archived    map[uuid.UUID][]Result
// }
//
// func (c *MemoryCatalog) ActiveJobById(id uuid.UUID) (Active, error) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	if _, found := c.activated[id]; !found {
// 		return Active{}, ErrNotFound
// 	}
//
// 	return c.activated[id], nil
// }
//
// func (c *MemoryCatalog) ActiveJobs() map[uuid.UUID]Active {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return c.activated
// }
//
// func (c *MemoryCatalog) ArchivedJobById(id uuid.UUID) ([]Result, error) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	if _, found := c.archived[id]; !found {
// 		return nil, ErrNotFound
// 	}
//
// 	return c.archived[id], nil
// }
//
// func (c *MemoryCatalog) ArchivedJobs() map[uuid.UUID][]Result {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return c.archived
// }
//
// func (c *MemoryCatalog) CountRegisteredJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.repository)
// }
//
// func (c *MemoryCatalog) CountSchedulableJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.schedulable)
// }
//
// func (c *MemoryCatalog) CountActivatedJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.activated)
// }
//
// func (c *MemoryCatalog) Delete(id uuid.UUID) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	if _, found := c.repository[id]; !found {
// 		return ErrNotFound
// 	}
//
// 	delete(c.repository, id)
// 	return nil
// }
//
// func (c *MemoryCatalog) Enable(id uuid.UUID) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	// var found bool
// 	if _, found := c.repository[id]; !found {
// 		return ErrNotFound
// 	}
// 	job := c.repository[id]
// 	job.Enable()
//
// 	c.repository[id] = job
// 	return nil
// }
//
// func (c *MemoryCatalog) RegisteredJobById(id uuid.UUID) (Definition, error) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	if _, found := c.repository[id]; !found {
// 		return Definition{}, ErrNotFound
// 	}
//
// 	return c.repository[id], nil
// }
//
// func (c *MemoryCatalog) RegisteredJobs() map[uuid.UUID]Definition {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return c.repository
// }
//
// func (c *MemoryCatalog) Register(job Definition) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	if _, found := c.repository[job.Uuid]; found {
// 		return ErrExist
// 	}
//
// 	c.repository[job.Uuid] = job
// 	if job.Enabled {
// 		c.schedulable[job.Uuid] = Schedulable{
// 			Uuid:     job.Uuid,
// 			Name:     job.Name,
// 			Schedule: job.Schedule,
// 			Tasks:    job.Tasks,
// 		}
// 	}
// 	return nil
// }
//
// func (c *MemoryCatalog) isRepeatable(id uuid.UUID) bool {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	job := c.repository[id]
// 	if job.Repeat {
// 		if job.MaxRuns != 0 && len(c.archived[id]) <= job.MaxRuns {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func (c *MemoryCatalog) schedule(id uuid.UUID) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	job := c.repository[id]
// 	if !job.Enabled {
// 		return
// 	}
//
// 	c.schedulable[id] = Schedulable{
// 		Uuid:     job.Uuid,
// 		Name:     job.Name,
// 		Schedule: job.Schedule,
// 		Tasks:    job.Tasks,
// 	}
// }

// REGISTERED JOB FUNCTIONS
// func (c *MemoryCatalog) CountRegisteredJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.registered)
// }

// func (c *MemoryCatalog) CountActiveJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.active)
// }
//
// func (c *MemoryCatalog) CountArchivedJobs() int {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return len(c.archive)
// }

// func (c *MemoryCatalog) Register(job Job) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	c.registered[job.Uuid] = job
// 	c.active[job.Uuid] = job
// }

// func (c *MemoryCatalog) Unregister(uuid uuid.UUID) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	delete(c.registered, uuid)
// }

// func (c *MemoryCatalog) GetArchivedJobs() []Job {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	return c.archive
// }
//
// // ACTIVE JOB FUNCTIONS
// func (c *MemoryCatalog) GetActiveJobs() []Job {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	var jobs = make([]Job, 0)
// 	for _, job := range c.active {
// 		jobs = append(jobs, job)
// 	}
// 	return jobs
// }

// func (c *MemoryCatalog) GetRegisteredJobs() []Job {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	var jobs = make([]Job, 0)
// 	for _, job := range c.registered {
// 		jobs = append(jobs, job)
// 	}
// 	return jobs
// }

// func (c *MemoryCatalog) UpdateActiveJob(job Job) {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()
//
// 	// fmt.Printf("Updating job \"%s\" - status \"%s\"\r\n", job.Name, job.Status)
// 	if job.Status == StatusCompleted || job.Status == StatusError {
// 		// Copy job to archive
// 		c.archive = append(c.archive, job)
//
// 		// Delete job from active jobs
// 		delete(c.active, job.Uuid)
//
// 		// Reactivate job if the job must be repeated, delete from repository if not
// 		if c.registered[job.Uuid].Repeat {
// 			// j := c.registered[job.Uuid]
// 			// task := PrintTask{Message: "##### " + strconv.Itoa(time.Now().Minute())}
// 			// j.Tasks.Tasks[1] = task
// 			// fmt.Printf("### Adding repeating job %s - status %s\r\n", j.Name, j.Status)
// 			// c.active[j.Uuid] = j
// 			c.active[job.Uuid] = c.registered[job.Uuid]
// 		} else {
// 			delete(c.registered, job.Uuid)
// 		}
// 	} else {
// 		c.active[job.Uuid] = job
// 	}
// }
