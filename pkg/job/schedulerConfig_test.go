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
	"testing"
	"time"
)

// TODO TestNewSchedulerConfig
func TestNewSchedulerConfig(t *testing.T) {
	c := NewSchedulerConfig()

	if c.MaxSchedulableJobs != 10 {
		t.Errorf("expected default scheduler config")
	}
}

func TestSchedulerConfig_GetNoJobsSchedulableDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("5s")

	if d != c.GetNoJobsSchedulableDelay() {
		t.Errorf("expected 5s duration, got %s", c.GetNoJobsSchedulableDelay().String())
	}
}

func TestSchedulerConfig_GetScheduleDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("1s")

	if d != c.GetScheduleDelay() {
		t.Errorf("expected 1s duration, got %s", c.GetScheduleDelay().String())
	}
}

// TODO TestSchedulerConfig_GetStartDelay
func TestSchedulerConfig_GetStartDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("0s")

	if d != c.GetStartDelay() {
		t.Errorf("expected 10s duration, got %s", c.GetStartDelay().String())
	}
}
