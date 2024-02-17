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

import "fmt"

func NewHandlerRepository() *HandlerRepository {
	return &HandlerRepository{
		handlerPool: make(map[string]*HandlerPool),
	}
}

type HandlerRepository struct {
	handlerPool map[string]*HandlerPool
}

func (r *HandlerRepository) GetHandlerNames() []string {
	keys := make([]string, len(r.handlerPool))
	for k := range r.handlerPool {
		keys = append(keys, k)
	}
	return keys
}

func (r *HandlerRepository) IsRegistered(handler string) bool {
	for _, h := range r.GetHandlerNames() {
		if h == handler {
			return true
		}
	}
	return false
}

func (r *HandlerRepository) RegisterHandlerPool(p *HandlerPool) {
	r.handlerPool[p.GetTaskType()] = p
}

func (r *HandlerRepository) Execute(t Task, pipeline chan *Pipeline) Task {
	handler, found := r.handlerPool[t.GetTaskType()]
	if !found {
		// If no handler is available, the program cannot continue
		panic(fmt.Sprintf("could not find handler for task %s", t.GetTaskType()))
	}
	return handler.Execute(t, pipeline)
}
