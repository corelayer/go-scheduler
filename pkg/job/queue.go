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
	"fmt"
	"sync"
)

func newQueue(ctx context.Context, max int) *queue {
	s := &queue{
		jobs: make([]Job, 0),
		chIn: make(chan Job, max),
	}
	go s.handleInput(ctx)
	return s
}

type queue struct {
	jobs []Job
	chIn chan Job

	mux sync.Mutex
}

func (q *queue) Length() int {
	return len(q.jobs)
}
func (q *queue) Capacity() int {
	return cap(q.jobs)
}

func (q *queue) Push(job Job) {
	q.chIn <- job
}

func (q *queue) Pop() (Job, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(q.jobs) > 0 {
		job := q.jobs[0]
		q.jobs = q.jobs[1:]
		return job, nil
	}
	return Job{}, fmt.Errorf("no jobs available")
}

func (q *queue) handleInput(ctx context.Context) {
	for {
		select {
		case job, ok := <-q.chIn:
			if !ok {
				return
			}
			q.mux.Lock()
			q.jobs = append(q.jobs, job)
			q.mux.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
