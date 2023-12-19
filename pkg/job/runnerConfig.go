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

func NewRunnerConfig(r *TaskHandlerRepository) RunnerConfig {
	c := RunnerConfig{
		maxConcurrentJobs:     10,
		idleDelayMilliseconds: 250,
		taskHandlerRepository: r,
	}
	return c
}

type RunnerConfig struct {
	maxConcurrentJobs     int
	idleDelayMilliseconds int
	taskHandlerRepository *TaskHandlerRepository
}

func (c RunnerConfig) WithMaxJobs(max int) RunnerConfig {
	c.maxConcurrentJobs = max
	return c
}

func (c RunnerConfig) WithIdleDelay(milliseconds int) RunnerConfig {
	c.idleDelayMilliseconds = milliseconds
	return c
}

func (c RunnerConfig) GetIdleDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.idleDelayMilliseconds) + "ms")
	return d
}
