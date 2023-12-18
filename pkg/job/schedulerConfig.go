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
		startDelayMilliseconds:    1,
		scheduleDelayMilliseconds: 1,
		idleDelayMilliseconds:     2,
		maxSchedulableJobs:        10,
	}
}

type SchedulerConfig struct {
	startDelayMilliseconds    int
	scheduleDelayMilliseconds int
	idleDelayMilliseconds     int
	maxSchedulableJobs        int
}

func (c SchedulerConfig) WithStartDelay(milliseconds int) SchedulerConfig {
	c.startDelayMilliseconds = milliseconds
	return c
}

func (c SchedulerConfig) WithScheduleDelay(milliseconds int) SchedulerConfig {
	c.scheduleDelayMilliseconds = milliseconds
	return c
}

func (c SchedulerConfig) WithIdleDelay(milliseconds int) SchedulerConfig {
	c.idleDelayMilliseconds = milliseconds
	return c
}

func (c SchedulerConfig) WithMaxJobs(max int) SchedulerConfig {
	c.maxSchedulableJobs = max
	return c
}

func (c SchedulerConfig) GetStartDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.startDelayMilliseconds) + "ms")
	return d
}

func (c SchedulerConfig) GetScheduleDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.scheduleDelayMilliseconds) + "ms")
	return d
}

func (c SchedulerConfig) GetIdleDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.idleDelayMilliseconds) + "ms")
	return d
}
