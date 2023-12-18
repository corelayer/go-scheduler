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
	"fmt"
	"reflect"
	"sync"
)

type Task interface {
}

type PrintTask struct {
	message     string
	readInput   bool
	writeOutput bool
}

type PrintTaskHandler struct{}

func (h PrintTaskHandler) Execute(t PrintTask, pipeline chan interface{}) Task {
	if t.readInput {
		select {
		case data := <-pipeline:
			fmt.Println(t.message)
			fmt.Println(data)
			if t.writeOutput {
				pipeline <- data
			}
		default:
			fmt.Println(t.message)
		}
	}
	return t
}

type EmptyTask struct {
	readInput   bool
	writeOutput bool
}

type EmptyTaskHandler struct{}

func (h EmptyTaskHandler) Execute(t EmptyTask, pipeline chan interface{}) Task {
	if t.readInput {
		select {
		case received := <-pipeline:
			if t.writeOutput {
				pipeline <- received
			}
		default:
		}
	}
	return t
}

type TaskHandlerRepository struct {
	handlerPool map[reflect.Type]*TaskHandlerPool
}

func (m *TaskHandlerRepository) RunTask(t Task, pipeline chan interface{}) Task {
	return m.handlerPool[reflect.TypeOf(t)].Execute(t, pipeline)
}

type TaskHandlerPool struct {
	handler         TaskHandler
	concurrentMax   int
	concurrentCount int
	mux             sync.Mutex
}

func (e *TaskHandlerPool) Execute(t Task, pipeline chan interface{}) Task {
	var output Task
	for {
		e.mux.Lock()
		if e.concurrentCount < e.concurrentMax {
			e.concurrentCount++
			e.mux.Unlock()
			output = e.handler.Execute(t, pipeline)
			break
		}
		e.mux.Unlock()
	}
	e.mux.Lock()
	e.concurrentCount--
	e.mux.Unlock()
	return output

}

type TaskHandler interface {
	Execute(t Task, pipeline chan interface{}) Task
}
