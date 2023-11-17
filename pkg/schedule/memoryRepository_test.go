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
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewRepository(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	r.mux.Lock()
	result := len(r.jobs)
	r.mux.Unlock()
	wanted := 0

	if result != wanted {
		t.Errorf("%s has capacity (%d), expected %d", t.Name(), result, wanted)
	}
}

func TestNewRepository2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	NewMemoryRepository(ctx)

	cancel()
}

func TestNewRepository3(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	close(r.chInput)
}

func TestNewRepository4(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	close(r.chDelete)
}

func TestNewRepository5(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	close(r.chUpdate)
}

func TestNewRepository6(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	close(r.chActivate)
}

func TestRepository_Add(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	r.Add(Job{
		Uuid:   uuid.New(),
		Name:   "test",
		Tasks:  nil,
		Status: JobStatusNone,
	})
}

func TestRepository_Delete(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	uuids := make([]uuid.UUID, 10)

	for i := 0; i < 10; i++ {
		id := uuid.New()
		uuids[i] = id
		r.Add(Job{
			Uuid: id,
			Name: strconv.Itoa(i),
		})
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(9)
	r.Delete(uuids[d])
}

func TestMemoryRepository_Schedulable(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	for i := 0; i < 10; i++ {
		r.addJob(Job{
			Uuid:   uuid.New(),
			Name:   strconv.Itoa(i),
			Tasks:  nil,
			Status: JobStatusSchedulable,
		})
	}
	result := r.Schedulable(0)
	wanted := 10

	if len(result) != wanted {
		t.Errorf("got %d schedulable jobs, expected %d", len(result), wanted)
	}
}

func TestMemoryRepository_Schedulable2(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	for i := 0; i < 10; i++ {
		r.addJob(Job{
			Uuid:   uuid.New(),
			Name:   strconv.Itoa(i),
			Tasks:  nil,
			Status: JobStatusSchedulable,
		})
	}

	result := r.Schedulable(5)
	wanted := 5

	if len(result) != wanted {
		t.Errorf("got %d schedulable jobs, expected %d", len(result), wanted)
	}
}

func TestRepository_Update(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	id := uuid.New()
	r.jobs[id] = Job{
		Uuid:   id,
		Name:   "test1",
		Tasks:  nil,
		Status: JobStatusNone,
	}

	r.Update(Job{
		Uuid:   id,
		Name:   "testUpdated",
		Tasks:  nil,
		Status: JobStatusCompleted,
	})

	r.mux.Lock()
	result := r.jobs[id]
	r.mux.Unlock()
	wanted := "testUpdated"

	if result.Name != wanted {
		t.Errorf("job name is %s, expected %s", result.Name, wanted)
	}
}

func TestRepository_Update2(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	uuid1 := uuid.New()
	r.jobs[uuid1] = Job{
		Uuid:   uuid1,
		Name:   "test1",
		Tasks:  nil,
		Status: JobStatusNone,
	}

	r.updateJob(Job{
		Uuid:   uuid.New(),
		Name:   "testUpdated",
		Tasks:  nil,
		Status: JobStatusCompleted,
	})

	r.mux.Lock()
	result := r.jobs[uuid1]
	r.mux.Unlock()
	wanted := "test1"

	if result.Name != wanted {
		t.Errorf("job name is %s, expected %s", result.Name, wanted)
	}
}

func TestMemoryRepository_deleteJob(t *testing.T) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	uuids := make([]uuid.UUID, 10)
	jobs := make(map[uuid.UUID]Job)
	for i := 0; i < 10; i++ {
		id := uuid.New()
		uuids[i] = id
		jobs[id] = Job{
			Uuid: id,
			Name: strconv.Itoa(i),
		}
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(9)

	r.deleteJob(uuids[d])

	r.mux.Lock()
	j := r.jobs
	r.mux.Unlock()

	stillExists := false
	for _, job := range j {
		if job.Uuid == uuids[d] {
			stillExists = true
			break
		}
	}

	if stillExists {
		t.Errorf("error deleting job %s", j[uuids[d]].Name)
	}
}

func BenchmarkRepository_Add(b *testing.B) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(Job{
			Uuid:   uuid.New(),
			Name:   "testJob",
			Tasks:  nil,
			Status: JobStatusNone,
		})
	}
}

func BenchmarkRepository_Update(b *testing.B) {
	ctx := context.Background()
	r := NewMemoryRepository(ctx)

	id := uuid.New()
	r.Add(Job{
		Uuid:   id,
		Name:   "testJob",
		Tasks:  nil,
		Status: JobStatusNone,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Update(Job{
			Uuid:   id,
			Name:   "testJob" + strconv.Itoa(i),
			Tasks:  nil,
			Status: 0,
		})
	}
}
