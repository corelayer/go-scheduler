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

func NewTaskHandlerRepository() *TaskHandlerRepository {
	return &TaskHandlerRepository{
		handlerPool: make(map[string]*TaskHandlerPool),
	}
}

type TaskHandlerRepository struct {
	handlerPool map[string]*TaskHandlerPool
}

func (r *TaskHandlerRepository) GetTaskHandlerNames() []string {
	keys := make([]string, len(r.handlerPool))
	for k := range r.handlerPool {
		keys = append(keys, k)
	}
	return keys
}

func (r *TaskHandlerRepository) IsRegistered(handler string) bool {
	handlers := r.GetTaskHandlerNames()
	for _, h := range handlers {
		if h == handler {
			return true
		}
	}
	return false
}

func (r *TaskHandlerRepository) RegisterTaskHandlerPool(p *TaskHandlerPool) {
	r.handlerPool[p.GetTaskType()] = p
}

func (r *TaskHandlerRepository) Execute(t Task, pipeline chan interface{}) Task {
	return r.handlerPool[t.GetTaskType()].Execute(t, pipeline)
}
