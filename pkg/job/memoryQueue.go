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
)

func NewMemoryQueue() *MemoryQueue {
	s := &MemoryQueue{
		jobs: make([]Job, 0),
	}
	return s
}

type MemoryQueue struct {
	jobs []Job
	mux  sync.Mutex
}

func (q *MemoryQueue) Length() int {
	q.mux.Lock()
	defer q.mux.Unlock()
	return len(q.jobs)

}
func (q *MemoryQueue) Capacity() int {
	q.mux.Lock()
	defer q.mux.Unlock()
	return cap(q.jobs)
}

func (q *MemoryQueue) Push(job Job) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.jobs = append(q.jobs, job)
}

func (q *MemoryQueue) Pop() (Job, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(q.jobs) > 0 {
		job := q.jobs[0]
		q.jobs = q.jobs[1:]
		return job, nil
	}
	return Job{}, fmt.Errorf("no jobs available")
}
