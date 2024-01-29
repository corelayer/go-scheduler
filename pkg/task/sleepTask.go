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
	"reflect"
	"strconv"
	"time"

	"github.com/corelayer/go-scheduler/pkg/job"
)

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

func (h SleepTaskHandler) Execute(t job.Task, pipeline chan interface{}) job.Task {
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
