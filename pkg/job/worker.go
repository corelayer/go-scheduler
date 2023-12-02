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
	"time"
)

func NewWorker(ctx context.Context, id int, chInput chan Job) *Worker {
	w := &Worker{
		id:      id,
		chInput: chInput,
	}
	go w.processJob(ctx)
	return w
}

type Worker struct {
	id      int
	chInput chan Job
}

func (w *Worker) processJob(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-w.chInput:
			if !ok {
				return
			}
			job.Status = StatusInProgress
			for _, t := range job.Tasks {
				t.Execute()
			}
		default:
			// fmt.Println("Waiting for jobs to be processed by worker", w.id)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
