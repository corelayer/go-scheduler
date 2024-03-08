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
		Tasks: tasks,
		mux:   &sync.Mutex{},
	}
}

type Sequence struct {
	Tasks     []Task
	executed  []Task
	active    bool
	activeIdx int
	mux       *sync.Mutex
}

func (s *Sequence) ActiveTask() Task {
	if s.IsActive() {
		s.mux.Lock()
		defer s.mux.Unlock()
		return s.Tasks[s.activeIdx]
	}
	return nil
}

func (s *Sequence) ActiveIndex() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.activeIdx
}

func (s *Sequence) All() []Task {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Tasks
}

func (s *Sequence) Count() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.Tasks)
}

func (s *Sequence) Execute(r *HandlerRepository, c *Intercom) {
	pipeline := make(chan *Pipeline, 1)
	defer close(pipeline)

	s.mux.Lock()
	s.active = true
	s.mux.Unlock()

	pipeline <- &Pipeline{Intercom: c}

	for i, t := range s.Tasks {
		s.mux.Lock()
		s.activeIdx = i
		s.mux.Unlock()

		result := r.Execute(t, pipeline)

		s.mux.Lock()
		s.executed = append(s.executed, result)
		s.mux.Unlock()
	}

	s.mux.Lock()
	s.active = false
	s.mux.Unlock()
}

func (s *Sequence) Executed() []Task {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.executed
}

func (s *Sequence) IsActive() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.active
}

func (s *Sequence) RegisterTask(t Task) {
	s.Tasks = append(s.Tasks, t)
}

func (s *Sequence) RegisterTasks(t []Task) {
	s.Tasks = append(s.Tasks, t...)
}

func (s *Sequence) ResetHistory() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.executed = nil
}
