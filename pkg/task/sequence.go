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
)

func NewSequence(tasks []Task) Sequence {
	return Sequence{
		tasks: tasks,
		mux:   &sync.Mutex{},
	}
}

type Sequence struct {
	tasks     []Task
	active    bool
	activeIdx int
	mux       *sync.Mutex
}

func (s Sequence) ActiveTask() Task {
	if s.IsActive() {
		s.mux.Lock()
		defer s.mux.Unlock()
		return s.tasks[s.activeIdx]
	}
	return nil
}

func (s Sequence) ActiveIndex() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.activeIdx
}

func (s Sequence) All() []Task {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.tasks
}

func (s Sequence) Count() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.tasks)
}

func (s Sequence) Execute(r *HandlerRepository, c *Intercom) {
	pipeline := make(chan *Pipeline, 1)
	defer close(pipeline)

	s.active = true
	pipeline <- &Pipeline{Intercom: c}

	for i, t := range s.tasks {
		s.mux.Lock()
		s.activeIdx = i
		s.tasks[i] = r.Execute(t, pipeline)
		s.mux.Unlock()
	}

	s.active = false
}

func (s Sequence) IsActive() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.active
}

func (s Sequence) RegisterTask(t Task) Sequence {
	s.tasks = append(s.tasks, t)
	return s
}

func (s Sequence) RegisterTasks(t []Task) Sequence {
	s.tasks = append(s.tasks, t...)
	return s
}
