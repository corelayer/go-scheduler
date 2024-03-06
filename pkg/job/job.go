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
	"github.com/corelayer/go-scheduler/pkg/task"
)

type Result struct {
	start    time.Time
	finish   time.Time
	runTime  time.Duration
	status   Status
	messages []task.Message
}

func NewJob(name string, s cron.Schedule, maxRuns int, tasks task.Sequence) Job {
	return Job{
		id:       uuid.New(),
		name:     name,
		enabled:  true,
		schedule: s,
		maxRuns:  maxRuns,
		status:   StatusInactive,
		results:  make([]Result, 0),
		tasks:    tasks,
	}
}

type Job struct {
	id       uuid.UUID
	name     string
	enabled  bool
	schedule cron.Schedule
	maxRuns  int
	status   Status
	results  []Result
	tasks    task.Sequence
}

func (j *Job) AddResult(r Result) {
	j.results = append(j.results, r)
}
func (j *Job) CountRuns() int {
	return len(j.results)
}

func (j *Job) CurrentResult() Result {
	return j.results[len(j.results)-1]
}

func (j *Job) Disable() {
	j.enabled = false
}

func (j *Job) Enable() {
	j.enabled = true
}

func (j *Job) IsAvailable() bool {
	return j.status == StatusAvailable
}

func (j *Job) IsEnabled() bool {
	return j.enabled
}

func (j *Job) IsInactive() bool {
	if j.maxRuns == 0 {
		return j.enabled && j.status == StatusInactive
	}

	if len(j.results) < j.maxRuns {
		return j.enabled && j.status == StatusInactive
	}

	return false
}

func (j *Job) IsRunnable() bool {
	return j.status == StatusRunnable
}

func (j *Job) IsSchedulable() bool {
	return j.status == StatusSchedulable && j.schedule.IsDue(time.Now())
}

func (j *Job) Results() []Result {
	return j.results
}

func (j *Job) UpdateResult(r Result) {
	j.results[len(j.results)-1] = r
}

func (j *Job) SetStatus(s Status) {
	j.status = s
}
