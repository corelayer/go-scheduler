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
)

func NewJob(name string, s cron.Schedule, maxRuns int) Job {
	return Job{
		id:       uuid.New(),
		name:     name,
		enabled:  true,
		schedule: s,
		maxRuns:  maxRuns,
	}
}

type Job struct {
	id       uuid.UUID
	name     string
	enabled  bool
	schedule cron.Schedule
	maxRuns  int

	scheduledTime  time.Time
	activationTime time.Time
	completedTime  time.Time
	runTime        time.Duration
}

func (j *Job) Activate() {
	j.activationTime = time.Now()
}

func (j *Job) Disable() {
	j.enabled = false
}

func (j *Job) Enable() {
	j.enabled = true
}
func (j *Job) Enabled() bool {
	return j.enabled
}

func (j *Job) IsRunnable() bool {
	if !j.enabled {
		return false
	}

	if j.schedule.IsDue(time.Now()) {
		return true
	}
	return false
}

func (j *Job) Finish() {
	j.completedTime = time.Now()
}

//
// import (
// 	"time"
//
// 	"github.com/google/uuid"
//
// 	"github.com/corelayer/go-scheduler/pkg/cron"
// 	"github.com/corelayer/go-scheduler/pkg/task"
// )
//
// func NewJob(id uuid.UUID, name string, enabled bool, schedule cron.Schedule, tasks []task.Task) Job {
// 	return Job{
// 		Uuid:     id,
// 		Name:     name,
// 		Enabled:  enabled,
// 		Status:   StatusNone,
// 		Schedule: schedule,
// 		Repeat:   false,
// 		Tasks:    task.NewSequence(tasks),
// 		Intercom: task.NewIntercom(),
// 	}
// }
//
// type Job struct {
// 	Uuid     uuid.UUID
// 	Name     string
// 	Enabled  bool
// 	Status   Status
// 	Schedule cron.Schedule
// 	Repeat   bool
// 	Tasks    task.Sequence
// 	Intercom *task.Intercom
// }
//
// func (j *Job) IsDue() bool {
// 	if !j.Enabled {
// 		return false
// 	}
//
// 	if j.Schedule.IsDue(time.Now()) {
// 		return true
// 	}
// 	return false
// }
//
// func (j *Job) IsPending() bool {
// 	return j.Status == StatusPending
// }
//
// func (j *Job) SetStatus(status Status) {
// 	j.Status = status
// }
//
// func (j *Job) SetTaskSequence(s task.Sequence) {
// 	j.Tasks = s
// }
//
// // TODO Definition rename to Job
// type Definition struct {
// 	Uuid     uuid.UUID
// 	Name     string
// 	Enabled  bool
// 	Status   Status
// 	Schedule cron.Schedule
// 	Repeat   bool
// 	MaxRuns  int // 0 = unlimited runs per schedule
// 	Tasks    task.Sequence
// }
//
// func (j *Definition) Disable() {
// 	j.Enabled = false
// }
//
// func (j *Definition) Enable() {
// 	j.Enabled = true
// }
//
// type Schedulable struct {
// 	Uuid     uuid.UUID
// 	Name     string
// 	Schedule cron.Schedule
// 	Tasks    task.Sequence
// }
//
// type Active struct {
// 	JobUuid  uuid.UUID
// 	RunUuid  uuid.UUID
// 	Name     string
// 	Status   Status
// 	Tasks    task.Sequence
// 	Intercom *task.Intercom
// }
//
// type Result struct {
// 	JobUuid  uuid.UUID
// 	RunUuid  uuid.UUID
// 	Name     string
// 	Status   Status
// 	Tasks    task.Sequence
// 	Intercom task.Intercom
// }
