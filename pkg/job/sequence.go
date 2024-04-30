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

	"github.com/corelayer/go-scheduler/pkg/task"
)

func NewSequence(tasks []task.Task) Sequence {
	return Sequence{
		Tasks: tasks,
		mux:   &sync.Mutex{},
	}
}

type Sequence struct {
	Tasks     []task.Task
	executed  []task.Task
	active    bool
	activeIdx int
	mux       *sync.Mutex
}

func (s *Sequence) ActiveTask() task.Task {
	if s.IsActive() {
		s.mux.Lock()
		defer s.mux.Unlock()
		return s.executed[s.activeIdx]
	}
	return nil
}

func (s *Sequence) ActiveIndex() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.activeIdx
}

func (s *Sequence) All() []task.Task {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.Tasks
}

func (s *Sequence) Count() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	return len(s.Tasks)
}

func (s *Sequence) CountExecuted() int {
	s.mux.Lock()
	defer s.mux.Unlock()

	executed := len(s.executed)
	if s.active && executed != 0 {
		executed--
	}
	return executed
}

func (s *Sequence) Executed() []task.Task {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.executed
}

func (s *Sequence) IsActive() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.active
}

func (s *Sequence) RegisterTask(t task.Task) {
	s.Tasks = append(s.Tasks, t)
}

func (s *Sequence) RegisterTasks(t []task.Task) {
	s.Tasks = append(s.Tasks, t...)
}

func (s *Sequence) ResetHistory() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.executed = nil
}
