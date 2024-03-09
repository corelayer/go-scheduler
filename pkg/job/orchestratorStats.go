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

import "github.com/google/uuid"

type GlobalStats struct {
	ConfiguredJobs  float64
	EnabledJobs     float64
	DisabledJobs    float64
	ActiveJobs      float64
	AvailableJobs   float64
	InactiveJobs    float64
	PendingJobs     float64
	RunnableJobs    float64
	SchedulableJobs float64
}

type TaskStats struct {
	Uuid      uuid.UUID
	Name      string
	Completed float64
	Total     float64
	HasErrors bool
}

type OrchestratorStats struct {
	Job   GlobalStats
	Tasks []TaskStats
}

func (o OrchestratorStats) HasTaskErrors() bool {
	for _, t := range o.Tasks {
		if t.HasErrors {
			return true
		}
	}
	return false
}

func (o OrchestratorStats) TasksCompleted() float64 {
	var tasksCompleted float64
	for _, t := range o.Tasks {
		tasksCompleted += t.Completed
	}
	return tasksCompleted
}

func (o OrchestratorStats) TasksTotal() float64 {
	var tasksTotal float64
	for _, t := range o.Tasks {
		tasksTotal += t.Total
	}
	return tasksTotal
}
