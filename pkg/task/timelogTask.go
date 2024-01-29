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
	"reflect"
	"time"

	"github.com/corelayer/go-scheduler/pkg/job"
)

type TimeLogTask struct {
	timestamp time.Time
}

func (t TimeLogTask) WriteToPipeline() bool {
	return true
}

func (t TimeLogTask) GetTaskType() string {
	return reflect.TypeOf(t).String()
}

type TimeLogTaskHandler struct{}

func (h TimeLogTaskHandler) Execute(t job.Task, pipeline chan interface{}) job.Task {
	task := t.(TimeLogTask)
	task.timestamp = time.Now()

	select {
	case data := <-pipeline:
		if task.WriteToPipeline() {
			pipeline <- data
		}
	default:
	}

	return task
}

func (h TimeLogTaskHandler) GetTaskType() string {
	return TimeLogTask{}.GetTaskType()
}
