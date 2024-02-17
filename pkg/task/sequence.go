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

package task

import (
	"sync"

	"github.com/corelayer/go-scheduler/pkg/status"
)

func NewSequence(tasks []Task) Sequence {
	return Sequence{
		Tasks: tasks,
		mux:   &sync.Mutex{},
	}
}

type Sequence struct {
	Tasks      []Task
	active     bool
	activeTask int
	mux        *sync.Mutex
}

func (s Sequence) GetActiveTaskIndex() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.activeTask
}

func (s Sequence) GetTaskStatus() []status.Status {
	s.mux.Lock()
	defer s.mux.Unlock()

	var output = make([]status.Status, len(s.Tasks))
	for i, t := range s.Tasks {
		output[i] = t.GetStatus()
	}
	return output
}

func (s Sequence) IsActive() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.active
}

func (s Sequence) RegisterTask(t Task) Sequence {
	s.Tasks = append(s.Tasks, t)
	return s
}

func (s Sequence) RegisterTasks(t []Task) Sequence {
	s.Tasks = append(s.Tasks, t...)
	return s
}

func (s Sequence) Run(r *HandlerRepository, c *Intercom) {
	pipeline := make(chan *Pipeline, 1)
	defer close(pipeline)

	s.active = true
	pipeline <- &Pipeline{Intercom: c}

	for i, t := range s.Tasks {
		s.mux.Lock()
		s.activeTask = i
		s.Tasks[i] = r.Execute(t, pipeline)
		s.mux.Unlock()
	}

	s.active = false
}
