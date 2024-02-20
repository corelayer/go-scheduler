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

func NewSchedulerConfig(maxJobs int, startupDelayMilliseconds int, chRunner chan Job, chUpdate chan Job) SchedulerConfig {
	return SchedulerConfig{
		ScheduleDelayMilliseconds: 500,
		StartupDelayMilliseconds:  startupDelayMilliseconds,
		MaxJobs:                   maxJobs,
		chRunner:                  chRunner,
		chUpdate:                  chUpdate,
	}
}

type SchedulerConfig struct {
	ScheduleDelayMilliseconds int
	StartupDelayMilliseconds  int
	MaxJobs                   int
	chRunner                  chan Job
	chUpdate                  chan Job
}

func (c SchedulerConfig) GetStartupDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.StartupDelayMilliseconds) + "ms")
	return d
}

func (c SchedulerConfig) GetScheduleDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.ScheduleDelayMilliseconds) + "ms")
	return d
}
