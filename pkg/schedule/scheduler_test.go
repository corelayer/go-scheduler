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

package schedule

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestNewScheduler(t *testing.T) {
	ctx := context.Background()
	c := SchedulerConfig{StartDelaySeconds: 5}
	r := NewMemoryRepository(ctx)

	for i := 0; i < 10; i++ {
		r.addJob(Job{
			Uuid:   uuid.New(),
			Name:   strconv.Itoa(i),
			Tasks:  nil,
			Status: JobStatusNone,
		})
		// t.Logf("adding job %d", i)
	}
	s := NewScheduler(ctx, c, r)

	result := len(s.jobs.All())
	wanted := 10

	if result != wanted {
		t.Errorf("%s has capacity (%d), expected %d", t.Name(), result, wanted)
	}
}

func TestNewScheduler2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := SchedulerConfig{StartDelaySeconds: 5}
	NewScheduler(ctx, c, nil)

	cancel()
}

//
// func TestNewScheduler3(t *testing.T) {
// 	ctx := context.Background()
// 	s := NewScheduler(ctx)
// 	close(s.chIn)
// }
//
// func TestScheduler_Add(t *testing.T) {
// 	ctx := context.Background()
// 	s := NewScheduler(ctx)
//
// 	for i := 0; i < 20; i++ {
// 		s.Add(&Job{
// 			Name:  "test",
// 			Tasks: []Task{TaskMock{}},
// 		})
// 	}
// }
//
// func BenchmarkScheduler_Add(b *testing.B) {
// 	ctx := context.Background()
// 	s := NewScheduler(ctx)
//
// 	for i := 0; i < b.N; i++ {
// 		s.Add(&Job{
// 			Name:  "test",
// 			Tasks: []Task{TaskMock{}},
// 		})
// 	}
//
// }

type TaskMock struct{}

func (m TaskMock) Execute() {
	return
}
func (m TaskMock) Notify(n chan JobStatus) {
	return
}
