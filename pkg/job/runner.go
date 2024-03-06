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
	"sync"
)

func NewRunner() *Runner {
	return &Runner{
		repository: NewRepository(),
		mux:        sync.Mutex{},
	}
}

type Runner struct {
	repository *Repository
	workers    []*Worker
	mux        sync.Mutex
}

func (r *Runner) Add(job Job) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.repository.Add(job)
}

func (r *Runner) Run(ctx context.Context) {

}

// func NewRunner(ctx context.Context, config RunnerConfig) (*Runner, error) {
// 	if config.Repository == nil {
// 		return nil, fmt.Errorf("invalid repository in RunnerConfig")
// 	}
//
// 	r := &Runner{
// 		config: config,
// 		// chQueue:         make(chan Job),
// 		// chWorkerUpdates: make(chan Job),
// 	}
//
// 	workers := make([]*Worker, config.MaxJobs)
// 	for i := 0; i < config.MaxJobs; i++ {
// 		// Create worker config
// 		wc, err := NewWorkerConfig(i, config.Repository)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Create worker to add to the pool
// 		workers[i], err = NewWorker(ctx, wc, r.config.chInput, r.config.chOutput)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	r.workers = workers
// 	return r, nil
//
// }
//
// type Runner struct {
// 	config  RunnerConfig
// 	workers []*Worker
// }
