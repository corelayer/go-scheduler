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
	"strconv"
	"time"
)

type Task interface {
	WriteToPipeline() bool
	GetTaskType() string
}

type TaskHandler interface {
	Execute(t Task, pipeline chan interface{}) Task
	GetTaskType() string
}

type SleepTask struct {
	Milliseconds int
	WriteOutput  bool
}

func (t SleepTask) WriteToPipeline() bool {
	return t.WriteOutput
}

func (t SleepTask) GetTaskType() string {
	return reflect.TypeOf(t).String()
}

type SleepTaskHandler struct{}

func (h SleepTaskHandler) GetTaskType() string {
	return SleepTask{}.GetTaskType()
}

func (h SleepTaskHandler) Execute(t Task, pipeline chan interface{}) Task {
	task := t.(SleepTask)
	d, _ := time.ParseDuration(strconv.Itoa(task.Milliseconds) + "ms")
	time.Sleep(d)

	select {
	case data := <-pipeline:
		if task.WriteToPipeline() {
			pipeline <- data
		}
	default:
	}

	return task
}

type PrintTask struct {
	Message     string
	ReadInput   bool
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

func (h PrintTaskHandler) Execute(t Task, pipeline chan interface{}) Task {
	task := t.(PrintTask)
	if task.ReadInput {
		select {
		case data := <-pipeline:
			fmt.Println(task.Message)
			fmt.Println(data)
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

type EmptyTask struct {
	readInput   bool
	writeOutput bool
}

func (t EmptyTask) WriteToPipeline() bool {
	return t.writeOutput
}

func (t EmptyTask) GetTaskType() string {
	return reflect.TypeOf(EmptyTask{}).String()
}

type EmptyTaskHandler struct{}

func (h EmptyTaskHandler) GetTaskType() string {
	return EmptyTask{}.GetTaskType()
}

func (h EmptyTaskHandler) Execute(t Task, pipeline chan interface{}) Task {
	task := t.(EmptyTask)
	if task.readInput {
		select {
		case received := <-pipeline:
			if task.WriteToPipeline() {
				pipeline <- received
			}
		default:
		}
	}
	return t
}
