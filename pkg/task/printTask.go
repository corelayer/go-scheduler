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

package task

import (
	"fmt"
	"reflect"

	"github.com/corelayer/go-scheduler/pkg/job"
)

type PrintTask struct {
	Message     string
	ReadInput   bool
	PrintInput  bool
	WriteOutput bool
}

func (t PrintTask) WriteToPipeline() bool {
	return t.WriteOutput
}

func (t PrintTask) GetTaskType() string {
	return reflect.TypeOf(t).String()
}

type PrintTaskHandler struct{}

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
