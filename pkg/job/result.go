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

package job

import (
	"time"

	"github.com/corelayer/go-scheduler/pkg/task"
)

type Result struct {
	Start    time.Time
	Finish   time.Time
	Status   Status
	Messages []task.Message
	Tasks    []task.Task
}

func (r Result) Runtime() time.Duration {
	return r.Finish.Sub(r.Start)
}
