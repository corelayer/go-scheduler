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

func TestNewSchedulerConfig(t *testing.T) {
	c := NewSchedulerConfig()

	if c.MaxSchedulableJobs != 10 {
		t.Errorf("expected default scheduler config")
	}
}

// TODO TestSchedulerConfig_WithMaxJobs
func TestSchedulerConfig_WithMaxJobs(t *testing.T) {

}

func TestSchedulerConfig_GetNoJobsSchedulableDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("2ms")

	if d != c.GetNoSchedulableJobsDelayDuration() {
		t.Errorf("expected 2ms duration, got %s", c.GetNoSchedulableJobsDelayDuration().String())
	}
}

func TestSchedulerConfig_GetScheduleDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("1ms")

	if d != c.GetScheduleDelayDuration() {
		t.Errorf("expected 1ms duration, got %s", c.GetScheduleDelayDuration().String())
	}
}

func TestSchedulerConfig_GetStartDelay(t *testing.T) {
	c := NewSchedulerConfig()
	d, _ := time.ParseDuration("1ms")

	if d != c.GetStartDelayDuration() {
		t.Errorf("expected 1ms duration, got %s", c.GetStartDelayDuration().String())
	}
}
