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

import "log/slog"

func NewTaskHandlerRepository() *TaskHandlerRepository {
	return &TaskHandlerRepository{
		handlerPool: make(map[string]*TaskHandlerPool),
	}
}

type TaskHandlerRepository struct {
	handlerPool map[string]*TaskHandlerPool
}

func (r *TaskHandlerRepository) RegisterTaskHandlerPool(p *TaskHandlerPool) {
	slog.Debug("register task handler pool", "type", p.GetTaskType())
	r.handlerPool[p.GetTaskType()] = p
}

func (r *TaskHandlerRepository) GetRegisteredHandlersNames() []string {
	keys := make([]string, len(r.handlerPool))
	for k, _ := range r.handlerPool {
		keys = append(keys, k)
	}
	return keys
}

func (r *TaskHandlerRepository) Execute(t Task, pipeline chan interface{}) Task {
	slog.Debug("handle task", "type", t.GetTaskType())
	_, found := r.handlerPool[t.GetTaskType()]
	if !found {
		slog.Error("could not find handler pool for task type", "task", t.GetTaskType(), "handlers", r.GetRegisteredHandlersNames())
		return t
	}
	return r.handlerPool[t.GetTaskType()].Execute(t, pipeline)
}
