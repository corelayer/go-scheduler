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
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestNewScheduler(t *testing.T) {
	ctx := context.Background()
	c := SchedulerConfig{startDelayMilliseconds: 5}
	r := NewMemoryCatalog()

	for i := 0; i < 10; i++ {
		r.Add(Job{
			Uuid:   uuid.New(),
			Name:   strconv.Itoa(i),
			Tasks:  TaskSequence{},
			Status: StatusNone,
		})
	}
	s, err := NewScheduler(ctx, c, r)
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}

	result := len(s.catalog.GetNotSchedulableJobs())
	wanted := 10

	if result != wanted {
		t.Errorf("%s has capacity (%d), expected %d", t.Name(), result, wanted)
	}
}

func TestNewScheduler2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := SchedulerConfig{startDelayMilliseconds: 5}
	NewScheduler(ctx, c, NewMemoryCatalog())

	cancel()
}

func TestNewScheduler3(t *testing.T) {
	ctx := context.Background()
	c := NewSchedulerConfig()
	s, err := NewScheduler(ctx, c, NewMemoryCatalog())
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}

	jobs := s.catalog.GetDueJobs(c.maxSchedulableJobs)

	result := len(jobs)
	wanted := 0
	if result != 0 {
		t.Errorf("found %d schedulable catalog, expected %d", result, wanted)
	}
}

func TestNewScheduler4(t *testing.T) {
	ctx := context.Background()

	r := NewMemoryCatalog()
	for i := 0; i < 10; i++ {
		r.Add(Job{
			Uuid:    uuid.New(),
			Name:    strconv.Itoa(i),
			Enabled: true,
			Tasks:   TaskSequence{},
			Status:  StatusIsDue,
		})
	}
	c := NewSchedulerConfig()
	s, err := NewScheduler(ctx, c, r)
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}

	jobs := s.catalog.GetDueJobs(c.maxSchedulableJobs)

	result := len(jobs)
	wanted := 10
	if result != wanted {
		t.Errorf("found %d schedulable catalog, expected %d", result, wanted)
	}

	jobs2 := s.catalog.GetDueJobs(c.maxSchedulableJobs)

	result2 := len(jobs2)
	wanted2 := 10
	if result2 != wanted2 {
		t.Errorf("found %d schedulable catalog, expected %d", result2, wanted2)
	}
}
