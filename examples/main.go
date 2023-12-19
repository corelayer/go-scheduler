package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
)

func createJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("@everysecond")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(250)
	tasks := []job.Task{
		job.PrintTask{
			Message:     fmt.Sprintf("Job %d - Task 1", i),
			ReadInput:   false,
			WriteOutput: true,
		},
		job.SleepTask{
			Milliseconds: d,
			WriteOutput:  true,
		},
		job.PrintTask{
			Message:     fmt.Sprintf("Job %d - Task 2", i),
			ReadInput:   true,
			WriteOutput: true,
		},
		job.SleepTask{
			Milliseconds: d,
			WriteOutput:  true,
		},
		job.PrintTask{
			Message:     fmt.Sprintf("Job %d - Task 3", i),
			ReadInput:   false,
			WriteOutput: false,
		},
	}

	return job.Job{
		Uuid:     id,
		Enabled:  true,
		Status:   job.StatusNone,
		Schedule: schedule,
		Name:     "Example Job 1",
		Tasks:    job.NewTaskSequence(tasks),
	}
}

func main() {

	c := job.NewMemoryCatalog()
	for i := 0; i < 100; i++ {
		c.Add(createJob(i))
	}

	p1 := job.NewTaskHandlerPool(job.PrintTaskHandler{}, 100)
	p2 := job.NewTaskHandlerPool(job.SleepTaskHandler{}, 100)

	r := job.NewTaskHandlerRepository()
	r.RegisterTaskHandlerPool(p1)
	r.RegisterTaskHandlerPool(p2)

	timeout, _ := time.ParseDuration("20s")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	config := job.OrchestratorConfig{MaxJobs: 20}
	_, err := job.NewOrchestrator(ctx, config, c, r)
	if err != nil {
		fmt.Println(err)
		cancel()
	}

	// time.Sleep(30 * time.Second)
	time.Sleep(10 * time.Second)
}
