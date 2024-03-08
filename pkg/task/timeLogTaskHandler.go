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
	"time"
)

const (
	TASKHANDLER_TIMELOG_MAX_CONCURRENT = 10000
)

func NewDefaultTimeLogTaskHandler() TimeLogTaskHandler {
	return TimeLogTaskHandler{
		maxConcurrency: TASKHANDLER_TIMELOG_MAX_CONCURRENT,
	}
}

func NewTimeLogTaskHandler(maxConcurrency int) TimeLogTaskHandler {
	return TimeLogTaskHandler{
		maxConcurrency: maxConcurrency,
	}
}

type TimeLogTaskHandler struct {
	maxConcurrency int
}

func (h TimeLogTaskHandler) Execute(t Task, p chan *Pipeline) Task {
	task := t.(TimeLogTask)
	pipeline := <-p
	task = h.processTask(task, pipeline)
	if task.WriteToPipeline() {
		p <- pipeline
	}
	return task.SetStatus(StatusCompleted)
}

func (h TimeLogTaskHandler) MaxConcurrent() int {
	return h.maxConcurrency
}

func (h TimeLogTaskHandler) Type() string {
	return TimeLogTask{}.Type()
}

func (h TimeLogTaskHandler) processTask(t TimeLogTask, p *Pipeline) TimeLogTask {
	timestamp := time.Now()

	p.Intercom.Add(Message{
		Message: "time",
		Type:    LogMessage,
		Task:    t.Name(),
		Data:    timestamp,
	})
	t.Timestamp = timestamp

	return t
}
