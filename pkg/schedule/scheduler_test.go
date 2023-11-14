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
	"testing"
)

func TestNewScheduler(t *testing.T) {
	ctx := context.Background()
	s := NewScheduler(ctx)

	result := cap(s.queue.jobs)
	wanted := 0

	if result != wanted {
		t.Errorf("%s has capacity (%d), expected %d", t.Name(), result, wanted)
	}
}

func TestNewScheduler2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	NewScheduler(ctx)

	cancel()
}

func TestNewScheduler3(t *testing.T) {
	ctx := context.Background()
	s := NewScheduler(ctx)
	close(s.chIn)
}

func TestScheduler_Add(t *testing.T) {
	ctx := context.Background()
	s := NewScheduler(ctx)

	for i := 0; i < 20; i++ {
		s.Add(&Job{
			Name:  "test",
			Tasks: []Task{TaskMock{}},
		})
	}
}

func BenchmarkScheduler_Add(b *testing.B) {
	ctx := context.Background()
	s := NewScheduler(ctx)

	for i := 0; i < b.N; i++ {
		s.Add(&Job{
			Name:  "test",
			Tasks: []Task{TaskMock{}},
		})
	}

}

type TaskMock struct{}

func (m TaskMock) Execute() {
	return
}
func (m TaskMock) Notify(n chan JobStatus) {
	return
}
