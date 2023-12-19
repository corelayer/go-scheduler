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
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
)

func TestNewOrchestrator(t *testing.T) {
	oc := OrchestratorConfig{
		MaxJobs:         10,
		SchedulerConfig: NewSchedulerConfig(),
		RunnerConfig:    NewRunnerConfig(&TaskHandlerRepository{}),
	}
	ctx, cancel := context.WithCancel(context.Background())

	c := NewMemoryCatalog()
	_, err := NewOrchestrator(ctx, oc, c)
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}

	schedule, _ := cron.NewSchedule("@everysecond")
	for i := 0; i < 100; i++ {
		c.Add(Job{
			Uuid:     uuid.New(),
			Enabled:  true,
			Status:   StatusNone,
			Schedule: schedule,
			Name:     strconv.Itoa(i),
			Tasks:    TaskSequence{},
		})
	}

	time.Sleep(15 * time.Second)
	cancel()
}
