/*
 * Copyright 2024 CoreLayer BV
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

type RepositoryReadWriter interface {
	RepositoryReader
	RepositoryWriter
}

type RepositoryReader interface {
	Get(id uuid.UUID) (Job, error)
}

type RepositoryWriter interface {
	Add(job Job) error
	Delete(id uuid.UUID) error
}

func NewRepository() *Repository {
	return &Repository{
		jobs: make([]Job, 0),
		mux:  sync.Mutex{},
	}
}

type Repository struct {
	jobs []Job
	mux  sync.Mutex
}

func (r *Repository) Add(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, j := range r.jobs {
		if job.Uuid == j.Uuid {
			return ErrExist
		}
	}

	r.jobs = append(r.jobs, job)
	return nil
}

func (r *Repository) Count() int {
	r.mux.Lock()
	defer r.mux.Unlock()

	return len(r.jobs)
}

func (r *Repository) Delete(id uuid.UUID) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	for i, j := range r.jobs {
		if j.Uuid == id {
			r.jobs = append(r.jobs[:i], r.jobs[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (r *Repository) Exists(id uuid.UUID) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, j := range r.jobs {
		if j.Uuid == id {
			return true
		}
	}

	return false
}

func (r *Repository) Get(id uuid.UUID) (Job, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, j := range r.jobs {
		if j.Uuid == id {
			return j, nil
		}
	}
	return Job{}, ErrNotFound
}

func (r *Repository) GetAll() []Job {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.jobs
}

func (r *Repository) Update(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	for i, j := range r.jobs {
		if j.Uuid == job.Uuid {
			r.jobs[i] = job
			return nil
		}
	}
	return ErrNotFound
}
