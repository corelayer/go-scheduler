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
)

func NewOrchestrator(ctx context.Context, c OrchestratorConfig, jobs CatalogReadWriter) (*Orchestrator, error) {
	var (
		err       error
		scheduler *Scheduler
		runner    *Runner
	)

	chRunner := make(chan Job, c.MaxJobs)
	chUpdate := make(chan Job, c.MaxJobs)

	sc := NewSchedulerConfig(c.MaxJobs, c.StartupDelayMilliseconds, chRunner, chUpdate)
	rc := NewRunnerConfig(c.MaxJobs, c.Repository, chRunner, chUpdate)

	scheduler, err = NewScheduler(ctx, sc, jobs)
	if err != nil {
		return nil, err
	}

	runner, err = NewRunner(ctx, rc)
	if err != nil {
		return nil, err
	}

	o := &Orchestrator{
		s:        scheduler,
		r:        runner,
		chRunner: chRunner,
		chUpdate: chUpdate,
	}

	return o, nil
}

type Orchestrator struct {
	s *Scheduler
	r *Runner

	chRunner chan Job
	chUpdate chan Job
}
