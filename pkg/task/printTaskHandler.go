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
)

const (
	MAX_CONCURRENT_TASKHANDLER_PRINT = 10000
)

func NewDefaultPrintTaskHandler() PrintTaskHandler {
	return PrintTaskHandler{
		maxConcurrent: MAX_CONCURRENT_TASKHANDLER_PRINT,
	}
}

func NewPrintTaskHandler(maxConcurrent int) PrintTaskHandler {
	return PrintTaskHandler{
		maxConcurrent: maxConcurrent,
	}
}

type PrintTaskHandler struct {
	maxConcurrent int
}

func (h PrintTaskHandler) MaxConcurrent() int {
	return h.maxConcurrent
}
func (h PrintTaskHandler) Type() string {
	return PrintTask{}.Type()
}

func (h PrintTaskHandler) Execute(t Task, p chan *Pipeline) Task {
	pipeline := <-p
	fmt.Println(t.(PrintTask).Message)

	if t.WriteToPipeline() {
		p <- pipeline
	}

	return t.SetStatus(StatusCompleted)
}
