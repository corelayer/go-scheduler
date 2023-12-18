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

func NewTaskSequence(tasks []Task) TaskSequence {
	return TaskSequence{
		tasks: tasks,
	}
}

type TaskSequence struct {
	pipeline chan interface{}
	tasks    []Task
}

func (s TaskSequence) RegisterTask(t Task) TaskSequence {
	s.tasks = append(s.tasks, t)
	return s
}

func (s TaskSequence) RegisterTasks(t []Task) TaskSequence {
	s.tasks = append(s.tasks, t...)
	return s
}

func (s TaskSequence) Run(manager *TaskHandlerRepository) {
	p := make(chan interface{})
	defer close(p)
	for i, t := range s.tasks {
		s.tasks[i] = manager.RunTask(t, p)
	}
}
