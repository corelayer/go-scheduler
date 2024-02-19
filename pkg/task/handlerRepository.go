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
	"fmt"
	"sync"
)

func NewHandlerRepository() *HandlerRepository {
	return &HandlerRepository{
		handlerPool: make(map[string]*HandlerPool),
		mux:         sync.Mutex{},
	}
}

type HandlerRepository struct {
	handlerPool map[string]*HandlerPool
	mux         sync.Mutex
}

func (r *HandlerRepository) Execute(t Task, pipeline chan *Pipeline) Task {
	r.mux.Lock()
	handler, found := r.handlerPool[t.Type()]
	r.mux.Unlock()

	if !found {
		// If no handler is available, the program cannot continue
		panic(fmt.Sprintf("could not find handler for task %s", t.Type()))
	}
	return handler.Execute(t, pipeline)
}

func (r *HandlerRepository) HandlerNames() []string {
	r.mux.Lock()
	defer r.mux.Unlock()

	keys := make([]string, len(r.handlerPool))
	for k := range r.handlerPool {
		keys = append(keys, k)
	}
	return keys
}

func (r *HandlerRepository) IsRegistered(handler string) bool {
	handlers := r.HandlerNames()

	for _, h := range handlers {
		if h == handler {
			return true
		}
	}
	return false
}

func (r *HandlerRepository) RegisterHandlerPool(p *HandlerPool) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, found := r.handlerPool[p.Type()]
	if !found {
		r.handlerPool[p.Type()] = p
		return nil
	}
	return fmt.Errorf("handlerpool %s is already registered", p.Type())
}

func (r *HandlerRepository) RegisterHandlerPools(pools []*HandlerPool) error {
	var err error
	for _, p := range pools {
		err = r.RegisterHandlerPool(p)
		if err != nil {
			return err
		}
	}
	return nil
}
