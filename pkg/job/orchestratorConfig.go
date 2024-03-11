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
	"fmt"
	"time"

	"github.com/corelayer/go-scheduler/pkg/task"
)

func NewOrchestratorConfig(maxJobs int, delay int, interval int, errF func(err error), msgF func(msg task.IntercomMessage)) (OrchestratorConfig, error) {
	var (
		err  error
		i, d time.Duration
	)
	if d, err = time.ParseDuration(fmt.Sprintf("%dms", delay)); err != nil {
		return OrchestratorConfig{}, err
	}
	if i, err = time.ParseDuration(fmt.Sprintf("%dms", interval)); err != nil {
		return OrchestratorConfig{}, err
	}

	return OrchestratorConfig{
		MaxJobs:          maxJobs,
		StartDelay:       d,
		ScheduleInterval: i,
		ErrorHandler:     errF,
		MessageHandler:   msgF,
	}, nil
}

type OrchestratorConfig struct {
	MaxJobs          int
	ScheduleInterval time.Duration // milliseconds
	StartDelay       time.Duration
	ErrorHandler     func(err error)
	MessageHandler   func(msg task.IntercomMessage)
}
