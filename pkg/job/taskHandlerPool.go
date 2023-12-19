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

import (
	"reflect"
	"sync"
)

func NewTaskHandlerPool(h TaskHandler, maxConcurrent int) *TaskHandlerPool {
	return &TaskHandlerPool{
		handler:         h,
		concurrentMax:   maxConcurrent,
		concurrentCount: 0,
		mux:             sync.Mutex{},
	}
}

type TaskHandlerPool struct {
	handler         TaskHandler
	concurrentMax   int
	concurrentCount int
	mux             sync.Mutex
}

func (p *TaskHandlerPool) GetTaskType() reflect.Type {
	return p.handler.GetTaskType()
}

func (p *TaskHandlerPool) Execute(t Task, pipeline chan interface{}) Task {
	var output Task
	for {
		p.mux.Lock()
		if p.concurrentCount < p.concurrentMax {
			p.concurrentCount++
			p.mux.Unlock()
			output = p.handler.Execute(t, pipeline)
			break
		}
		p.mux.Unlock()
	}
	p.mux.Lock()
	p.concurrentCount--
	p.mux.Unlock()
	return output
}
