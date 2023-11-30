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
	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
)

type Job struct {
	Uuid     uuid.UUID
	Enabled  bool
	Status   Status
	Schedule cron.Schedule
	Name     string
	Tasks    []TaskRunner
}

func (j *Job) Activate() bool {
	return true
}

func (j *Job) IsSchedulable() bool {
	if !j.Enabled {
		return false
	}
	if j.Status == StatusSchedulable {
		return true
	}
	return false
}
