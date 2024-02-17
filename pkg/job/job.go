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
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/status"
	"github.com/corelayer/go-scheduler/pkg/task"
)

type Job struct {
	Uuid     uuid.UUID
	Enabled  bool
	Status   status.Status
	Schedule cron.Schedule
	Repeat   bool
	Name     string
	Tasks    task.Sequence
	Intercom *task.Intercom
}

func (j *Job) IsDue() bool {
	if !j.Enabled {
		return false
	}

	if j.Schedule.IsDue(time.Now()) {
		return true
	}
	return false
}

func (j *Job) SetStatus(status status.Status) {
	j.Status = status
}

func (j *Job) IsPending() bool {
	return j.Status == status.StatusPending
}
