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
	"strconv"
	"time"
)

const (
	MaxConcurrentTaskHandlerSleep = 10000
)

func NewDefaultSleepTaskHandler() SleepTaskHandler {
	return SleepTaskHandler{
		maxConcurrent: MaxConcurrentTaskHandlerSleep,
	}
}

func NewSleepTaskHandler(maxConcurrent int) SleepTaskHandler {
	return SleepTaskHandler{
		maxConcurrent: maxConcurrent,
	}
}

type SleepTaskHandler struct {
	maxConcurrent int
}

func (h SleepTaskHandler) Execute(t Task, p chan *Pipeline) Task {
	d, _ := time.ParseDuration(strconv.Itoa(t.(SleepTask).Milliseconds) + "ms")
	time.Sleep(d)

	pipeline := <-p
	if t.WriteToPipeline() {
		p <- pipeline
	}

	return t.SetStatus(StatusCompleted)
}

func (h SleepTaskHandler) MaxConcurrent() int {
	return h.maxConcurrent
}

func (h SleepTaskHandler) Type() string {
	return SleepTask{}.Type()
}
