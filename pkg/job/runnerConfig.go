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

func NewRunnerConfig() RunnerConfig {
	c := RunnerConfig{
		MaxConcurrentJobs:               10,
		NoRunnableJobsDelayMilliseconds: 250,
	}
	return c
}

type RunnerConfig struct {
	MaxConcurrentJobs               int
	NoRunnableJobsDelayMilliseconds int
}

func (c RunnerConfig) WithMaxJobs(max int) RunnerConfig {
	c.MaxConcurrentJobs = max
	return c
}

func (c RunnerConfig) GetNoRunnableJobsDelayDuration() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.NoRunnableJobsDelayMilliseconds) + "ms")
	return d
}
