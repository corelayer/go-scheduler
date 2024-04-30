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
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
)

func NewJob(name string, s cron.Schedule, maxRuns int, tasks Sequence) Job {
	return Job{
		Uuid:     uuid.New(),
		Name:     name,
		Enabled:  true,
		Schedule: s,
		MaxRuns:  maxRuns,
		Status:   StatusInactive,
		Tasks:    tasks,
		History:  make([]Result, 0),
		mux:      &sync.Mutex{},
	}
}

type Job struct {
	Uuid     uuid.UUID
	Name     string
	Enabled  bool
	Schedule cron.Schedule
	MaxRuns  int
	Status   Status
	Tasks    Sequence
	History  []Result
	mux      *sync.Mutex
}

func (j *Job) AddResult(r Result) {
	j.mux.Lock()
	defer j.mux.Unlock()

	j.History = append(j.History, r)
}
func (j *Job) CountRuns() int {
	j.mux.Lock()
	defer j.mux.Unlock()

	return len(j.History)
}

func (j *Job) CurrentResult() Result {
	j.mux.Lock()
	defer j.mux.Unlock()

	if len(j.History) > 0 {
		return j.History[len(j.History)-1]
	}
	return Result{}
}

func (j *Job) Disable() {
	j.mux.Lock()
	defer j.mux.Unlock()

	j.Enabled = false
}

func (j *Job) Enable() {
	j.mux.Lock()
	defer j.mux.Unlock()

	j.Enabled = true
}

func (j *Job) IsActive() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.Status == StatusActive
}

func (j *Job) IsAvailable() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.Status == StatusAvailable
}

func (j *Job) IsEligible() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	if j.MaxRuns == 0 {
		return j.Enabled
	}

	if len(j.History) < j.MaxRuns {
		return j.Enabled
	}
	return false
}

func (j *Job) IsEnabled() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.Enabled
}

func (j *Job) IsInactive() bool {
	eligible := j.IsEligible()

	j.mux.Lock()
	defer j.mux.Unlock()
	return eligible && j.Status == StatusInactive
}

func (j *Job) IsPending() bool {
	j.mux.Lock()
	defer j.mux.Lock()

	return j.Status == StatusPending
}

func (j *Job) IsRunnable() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.Status == StatusRunnable
}

func (j *Job) IsSchedulable() bool {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.Status == StatusSchedulable && j.Schedule.IsDue(time.Now())
}

func (j *Job) AllResults() []Result {
	j.mux.Lock()
	defer j.mux.Unlock()

	return j.History
}

func (j *Job) UpdateResult(r Result) {
	j.mux.Lock()
	defer j.mux.Unlock()

	if len(j.History) > 0 {
		j.History[len(j.History)-1] = r
	} else {
		j.History = append(j.History, r)
	}
}

func (j *Job) SetStatus(s Status) {
	j.mux.Lock()
	defer j.mux.Unlock()

	j.Status = s
}
