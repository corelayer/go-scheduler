// /*
//   - Copyright 2023 CoreLayer BV
//     *
//   - Licensed under the Apache License, Version 2.0 (the "License");
//   - you may not use this file except in compliance with the License.
//   - You may obtain a copy of the License at
//     *
//   - http://www.apache.org/licenses/LICENSE-2.0
//     *
//   - Unless required by applicable law or agreed to in writing, software
//   - distributed under the License is distributed on an "AS IS" BASIS,
//   - WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   - See the License for the specific language governing permissions and
//   - limitations under the License.
//     */
package job

//
// func TestNewMemoryCatalog(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	r.mux.Lock()
// 	result := len(r.registered)
// 	r.mux.Unlock()
// 	wanted := 0
//
// 	if result != wanted {
// 		t.Errorf("%s has capacity (%d), expected %d", t.Name(), result, wanted)
// 	}
// }
//
// func TestMemoryCatalog_Add(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	r.Register(Job{
// 		Uuid:   uuid.New(),
// 		Name:   "test",
// 		Tasks:  task.Sequence{},
// 		Status: StatusNone,
// 	})
// }
//
// func TestMemoryCatalog_Count(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	uuids := make([]uuid.UUID, 10)
//
// 	for i := 0; i < 10; i++ {
// 		id := uuid.New()
// 		uuids[i] = id
// 		r.Register(Job{
// 			Uuid: id,
// 			Name: strconv.Itoa(i),
// 		})
// 	}
//
// 	result := r.CountRegisteredJobs()
// 	wanted := 10
//
// 	if result != wanted {
// 		t.Errorf("got %d schedulable jobs, expected %d", result, wanted)
// 	}
//
// }
//
// func TestMemoryCatalog_Delete(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	uuids := make([]uuid.UUID, 10)
//
// 	for i := 0; i < 10; i++ {
// 		id := uuid.New()
// 		uuids[i] = id
// 		r.Register(Job{
// 			Uuid: id,
// 			Name: strconv.Itoa(i),
// 		})
// 	}
//
// 	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	d := rnd.Intn(9)
// 	r.Unregister(uuids[d])
// }
//
// func TestMemoryCatalog_Update(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	id := uuid.New()
// 	r.Register(Job{
// 		Uuid:   id,
// 		Name:   "test1",
// 		Tasks:  task.Sequence{},
// 		Status: StatusNone,
// 	})
//
// 	r.UpdateActiveJob(Job{
// 		Uuid:   id,
// 		Name:   "testUpdated",
// 		Tasks:  task.Sequence{},
// 		Status: StatusPending,
// 	})
//
// 	r.mux.Lock()
// 	result := r.active[id]
// 	r.mux.Unlock()
// 	wanted := "testUpdated"
//
// 	if result.Name != wanted {
// 		t.Errorf("job name is %s, expected %s", result.Name, wanted)
// 	}
// }
//
// func TestMemoryCatalog_Update2(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	uuid1 := uuid.New()
// 	r.Register(Job{
// 		Uuid:   uuid1,
// 		Name:   "test1",
// 		Tasks:  task.Sequence{},
// 		Status: StatusNone,
// 	})
//
// 	r.UpdateActiveJob(Job{
// 		Uuid:   uuid.New(),
// 		Name:   "testUpdated",
// 		Tasks:  task.Sequence{},
// 		Status: StatusCompleted,
// 	})
//
// 	r.mux.Lock()
// 	result := r.active[uuid1]
// 	r.mux.Unlock()
// 	wanted := "test1"
//
// 	if result.Name != wanted {
// 		t.Errorf("job name is %s, expected %s", result.Name, wanted)
// 	}
// }
//
// func TestMemoryCatalog_deleteJob(t *testing.T) {
// 	r := NewMemoryCatalog()
//
// 	uuids := make([]uuid.UUID, 10)
// 	jobs := make(map[uuid.UUID]Job)
// 	for i := 0; i < 10; i++ {
// 		id := uuid.New()
// 		uuids[i] = id
// 		jobs[id] = Job{
// 			Uuid: id,
// 			Name: strconv.Itoa(i),
// 		}
// 	}
//
// 	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	d := rnd.Intn(9)
//
// 	r.Unregister(uuids[d])
//
// 	r.mux.Lock()
// 	j := r.registered
// 	r.mux.Unlock()
//
// 	stillExists := false
// 	for _, job := range j {
// 		if job.Uuid == uuids[d] {
// 			stillExists = true
// 			break
// 		}
// 	}
//
// 	if stillExists {
// 		t.Errorf("error deleting job %s", j[uuids[d]].Name)
// 	}
// }
//
// func BenchmarkMemoryCatalog_Add(b *testing.B) {
// 	r := NewMemoryCatalog()
// 	s, _ := cron.NewSchedule("@everysecond")
//
// 	var id []uuid.UUID
// 	for i := 0; i < b.N; i++ {
// 		id = append(id, uuid.New())
// 	}
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		r.Register(Job{
// 			Uuid:     id[i],
// 			Name:     "testJob",
// 			Tasks:    task.Sequence{},
// 			Status:   StatusNone,
// 			Schedule: s,
// 			Enabled:  false,
// 		})
// 	}
// }
//
// func BenchmarkMemoryCatalog_Update(b *testing.B) {
// 	r := NewMemoryCatalog()
// 	s, _ := cron.NewSchedule("@everysecond")
//
// 	id := uuid.New()
// 	r.Register(Job{
// 		Uuid:     id,
// 		Name:     "testJob",
// 		Tasks:    task.Sequence{},
// 		Status:   StatusNone,
// 		Schedule: s,
// 		Enabled:  false,
// 	})
//
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		r.UpdateActiveJob(Job{
// 			Uuid:   id,
// 			Name:   "a",
// 			Tasks:  task.Sequence{},
// 			Status: StatusPending,
// 		})
// 	}
// }
