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

import "github.com/corelayer/go-scheduler/pkg/job"

const (
	EMPTY_TASKHANDLER_CONCURRENCY = 10000
)

func NewDefaultEmptyTaskHandler() EmptyTaskHandler {
	return EmptyTaskHandler{
		maxConcurrency: EMPTY_TASKHANDLER_CONCURRENCY,
	}
}

func NewEmptyTaskHandler(maxConcurrency int) EmptyTaskHandler {
	return EmptyTaskHandler{
		maxConcurrency: maxConcurrency,
	}
}

type EmptyTaskHandler struct {
	maxConcurrency int
}

func (h EmptyTaskHandler) GetMaxConcurrency() int {
	return h.maxConcurrency
}
func (h EmptyTaskHandler) GetTaskType() string {
	return EmptyTask{}.GetTaskType()
}

func (h EmptyTaskHandler) Execute(t job.Task, pipeline chan interface{}) job.Task {
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