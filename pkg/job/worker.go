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
	"fmt"
	"strconv"
	"time"

	"github.com/corelayer/go-scheduler/pkg/task"
)

func NewWorkerConfig(id int, r *task.HandlerRepository) (WorkerConfig, error) {
	if r == nil {
		return WorkerConfig{}, fmt.Errorf("invalid task.HandlerRepository")
	}
	return WorkerConfig{
		id:                        id,
		repository:                r,
		idleSleepTimeMilliseconds: 250,
	}, nil
}

type WorkerConfig struct {
	id                        int
	repository                *task.HandlerRepository
	idleSleepTimeMilliseconds int
}

func (c *WorkerConfig) GetIdleDelay() time.Duration {
	d, _ := time.ParseDuration(strconv.Itoa(c.idleSleepTimeMilliseconds) + "ms")
	return d
}

func NewWorker(ctx context.Context, config WorkerConfig, chInput chan Job, chUpdate chan Job) (*Worker, error) {
	w := &Worker{
		Config:   config,
		chInput:  chInput,
		chUpdate: chUpdate,
	}

	if config.repository == nil {
		return nil, fmt.Errorf("invalid repository in WorkerConfig")
	}
	go w.run(ctx)
	return w, nil
}

type Worker struct {
	Config   WorkerConfig
	chInput  chan Job
	chUpdate chan Job
}

func (w *Worker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-w.chInput:
			if !ok {
				return
			}

			if job.Intercom == nil {
				job.Intercom = task.NewIntercom()
			}

			// Run all task for job
			job.Tasks.Execute(w.Config.repository, job.Intercom)

			if job.Intercom.HasErrors() {
				job.SetStatus(StatusError)
			} else {
				job.SetStatus(StatusCompleted)
			}
			w.chUpdate <- job
		}
	}
}
