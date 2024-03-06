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
		data: make(map[uuid.UUID]Job),
		mux:  sync.Mutex{},
	}
}

type Repository struct {
	data map[uuid.UUID]Job
	mux  sync.Mutex
}

func (r *Repository) Add(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, found := r.data[job.id]; found {
		return ErrExist
	}

	r.data[job.id] = job
	return nil
}

func (r *Repository) Count() int {
	r.mux.Lock()
	defer r.mux.Unlock()

	return len(r.data)
}

func (r *Repository) Delete(id uuid.UUID) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, found := r.data[id]; found {
		return ErrNotFound
	}

	delete(r.data, id)
	return nil
}

func (r *Repository) Exists(id uuid.UUID) bool {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, found := r.data[id]; !found {
		return false
	}

	return true
}

func (r *Repository) Get(id uuid.UUID) (Job, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, found := r.data[id]; !found {
		return Job{}, ErrNotFound
	}

	return r.data[id], nil
}

func (r *Repository) GetAll() map[uuid.UUID]Job {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.data
}

func (r *Repository) Update(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, found := r.data[job.id]; !found {
		return ErrNotFound
	}

	r.data[job.id] = job
	return nil
}
