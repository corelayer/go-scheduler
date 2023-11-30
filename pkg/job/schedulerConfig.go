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
	"strconv"
	"time"
)

func NewSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		StartDelaySeconds:             0,
		ScheduleDelaySeconds:          1,
		NoSchedulableJobsDelaySeconds: 5,
		MaxSchedulableJobs:            10,
	}
}

type SchedulerConfig struct {
	StartDelaySeconds             int
	ScheduleDelaySeconds          int
	NoSchedulableJobsDelaySeconds int
	MaxSchedulableJobs            int
}

func (c *SchedulerConfig) GetStartDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.StartDelaySeconds) + "s")
	return d
}

func (c *SchedulerConfig) GetScheduleDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.ScheduleDelaySeconds) + "s")
	return d
}

func (c *SchedulerConfig) GetNoJobsSchedulableDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.NoSchedulableJobsDelaySeconds) + "s")
	return d
}
