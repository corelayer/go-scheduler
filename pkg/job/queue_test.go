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
	"testing"
)

func TestQueue_Add(t *testing.T) {
	ctx := context.Background()
	jl := newQueue(ctx, 10)

	for i := 0; i < 20; i++ {
		jl.Push(Job{Name: "test"})
	}
}

func TestNewQueue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	newQueue(ctx, 10)
	cancel()
}

func TestNewQueue2(t *testing.T) {
	ctx := context.Background()
	jl := newQueue(ctx, 10)
	close(jl.chIn)
}

func TestQueue_Get(t *testing.T) {
	ctx := context.Background()
	jl := newQueue(ctx, 10)
	tests := make([]string, 10)

	for i := 0; i < 10; i++ {
		jl.Push(Job{Name: "test"})
		tests[i] = "test"
	}

	tests = append(tests, "job_unknown")

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			result, err := jl.Pop()
			if err != nil {
				return
			}

			switch out := interface{}(result).(type) {
			case Job:
				return
			default:
				t.Errorf("got %s, expected Job", out)
			}

		})
	}
}

func TestQueue_Length(t *testing.T) {
	ctx := context.Background()
	jl := newQueue(ctx, 10)

	result := jl.Length()
	wanted := 0

	if result != wanted {
		t.Errorf("length = %d, expected %d", result, wanted)
	}
}

func TestQueue_Capacity(t *testing.T) {
	ctx := context.Background()
	jl := newQueue(ctx, 10)

	result := jl.Capacity()
	wanted := 0

	if result != wanted {
		t.Errorf("capacity = %d, expected %d", result, wanted)
	}
}

// func BenchmarkJobQueue_Add(b *testing.B) {
// 	ctx := context.Background()
// 	jl := newQueue(ctx)
//
// 	for i := 0; i < b.N; i++ {
// 		jl.Push(Job{Name: "test"})
// 	}
// }
