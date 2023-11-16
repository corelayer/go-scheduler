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

package schedule

import "github.com/google/uuid"

type Job struct {
	Uuid   uuid.UUID
	Name   string
	Tasks  []Task
	Status JobStatus
}

func (j *Job) IsSchedulable() bool {
	if j.Status == JobStatusSchedulable {
		return true
	}
	return false
}

type Task interface {
	Execute()
	Notify(n chan JobStatus)
}

type JobStatus int

const (
	JobStatusNone JobStatus = iota
	JobStatusSchedulable
	JobStatusScheduled
	JobStatusPending
	JobStatusStarted
	JobStatusInProgress
	JobStatusCompleted
)
