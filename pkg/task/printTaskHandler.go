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

	"github.com/corelayer/go-scheduler/pkg/job"
)

const (
	PRINT_TASKHANDLER_MAX_CONCURRENCY = 10000
)

func NewDefaultPrintTaskHandler() PrintTaskHandler {
	return PrintTaskHandler{
		maxConcurrency: PRINT_TASKHANDLER_MAX_CONCURRENCY,
	}
}

func NewPrintTaskHandler(maxConcurrency int) PrintTaskHandler {
	return PrintTaskHandler{
		maxConcurrency: maxConcurrency,
	}
}

type PrintTaskHandler struct {
	maxConcurrency int
}

func (h PrintTaskHandler) GetMaxConcurrency() int {
	return h.maxConcurrency
}
func (h PrintTaskHandler) GetTaskType() string {
	return PrintTask{}.GetTaskType()
}

func (h PrintTaskHandler) Execute(t job.Task, pipeline chan interface{}) job.Task {
	task := t.(PrintTask)
	if task.ReadInput {
		select {
		case data := <-pipeline:
			fmt.Println(task.Message)
			if task.PrintInput {
				fmt.Println(data)
			}
			if task.WriteToPipeline() {
				pipeline <- data
			}
		default:
			fmt.Println(task.Message)
		}
	} else {
		fmt.Println(task.Message)
	}
	return task
}
