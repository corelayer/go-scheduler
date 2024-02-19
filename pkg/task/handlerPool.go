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

func NewHandlerPool(h Handler) *HandlerPool {
	return &HandlerPool{
		handler:         h,
		concurrentMax:   h.MaxConcurrent(),
		concurrentCount: 0,
		mux:             &sync.Mutex{},
	}
}

type HandlerPool struct {
	handler         Handler
	concurrentMax   int
	concurrentCount int
	mux             *sync.Mutex
}

func (p *HandlerPool) ActiveHandlers() int {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.concurrentCount
}

func (p *HandlerPool) AvailableHandlers() int {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.concurrentMax - p.concurrentCount
}

func (p *HandlerPool) Execute(t Task, pipeline chan *Pipeline) Task {
	for {
		p.mux.Lock()
		if p.concurrentCount < p.concurrentMax {
			p.concurrentCount++
			p.mux.Unlock()
			t = p.handler.Execute(t, pipeline)
			break
		}
		p.mux.Unlock()
	}
	p.mux.Lock()
	p.concurrentCount--
	p.mux.Unlock()
	return t
}

func (p *HandlerPool) Type() string {
	return p.handler.Type()
}
