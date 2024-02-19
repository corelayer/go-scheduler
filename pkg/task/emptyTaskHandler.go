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

import "github.com/corelayer/go-scheduler/pkg/status"

const (
	MAX_CONCURRENT_TASKHANDLER_EMPTY = 10000
)

func NewDefaultEmptyTaskHandler() EmptyTaskHandler {
	return EmptyTaskHandler{
		maxConcurrent: MAX_CONCURRENT_TASKHANDLER_EMPTY,
	}
}

func NewEmptyTaskHandler(maxConcurrent int) EmptyTaskHandler {
	return EmptyTaskHandler{
		maxConcurrent: maxConcurrent,
	}
}

type EmptyTaskHandler struct {
	maxConcurrent int
}

func (h EmptyTaskHandler) Execute(t Task, p chan *Pipeline) Task {
	select {
	case pipeline := <-p:
		if t.WriteToPipeline() {
			p <- pipeline
		}
	default:
	}
	return t.SetStatus(status.StatusCompleted)
}

func (h EmptyTaskHandler) MaxConcurrent() int {
	return h.maxConcurrent
}
func (h EmptyTaskHandler) Type() string {
	return EmptyTask{}.Type()
}
