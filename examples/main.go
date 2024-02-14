package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
	"github.com/corelayer/go-scheduler/pkg/task"
)

func createJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("@everysecond")
	// rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// d := rnd.Intn(250)
	tasks := []job.Task{
		task.TimeLogTask{},
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 1", i),
		// 	ReadInput:   false,
		// 	WriteOutput: false,
		// },
		task.EmptyTask{},
		// task.SleepTask{
		// 	Milliseconds: d,
		// 	WriteOutput:  false,
		// },
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 2", i),
		// 	ReadInput:   true,
		// 	WriteOutput: true,
		// },
		// task.SleepTask{
		// 	Milliseconds: d,
		// 	WriteOutput:  true,
		// },
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 3", i),
		// 	ReadInput:   false,
		// 	WriteOutput: false,
		// },
		task.TimeLogTask{},
	}

	return job.Job{
		Uuid:     id,
		Enabled:  true,
		Status:   job.StatusNone,
		Schedule: schedule,
		Repeat:   false,
		Name:     "Example Job " + strconv.Itoa(i),
		Tasks:    job.NewTaskSequence(tasks),
	}
}

func createRepeatableJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("* * * * * *")
	// rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// d := rnd.Intn(1000)
	tasks := []job.Task{
		task.TimeLogTask{},
		task.PrintTask{
			Message:     fmt.Sprintf("### Repeatable job %d - Print Task", i),
			ReadInput:   false,
			WriteOutput: true,
			PrintInput:  true,
		},
		// task.SleepTask{
		// 	Milliseconds: d,
		// 	WriteOutput:  true,
		// },
		task.TimeLogTask{},
	}

	return job.Job{
		Uuid:     id,
		Enabled:  true,
		Status:   job.StatusNone,
		Schedule: schedule,
		Repeat:   true,
		Name:     "### Repeatable job " + strconv.Itoa(i),
		Tasks:    job.NewTaskSequence(tasks),
	}
}

func main() {

	c := job.NewMemoryCatalog()
	for i := 0; i < 1000000; i++ {
		c.Register(createJob(i))
	}
	// c.Register(createRepeatableJob(1))

	p1 := job.NewTaskHandlerPool(task.NewDefaultEmptyTaskHandler())
	p2 := job.NewTaskHandlerPool(task.NewDefaultSleepTaskHandler())
	// p3 := job.NewTaskHandlerPool(task.NewDefaultTimeLogTaskHandler())
	p4 := job.NewTaskHandlerPool(task.NewDefaultEmptyTaskHandler())

	r := job.NewTaskHandlerRepository()
	r.RegisterTaskHandlerPool(p1)
	r.RegisterTaskHandlerPool(p2)
	// r.RegisterTaskHandlerPool(p3)
	r.RegisterTaskHandlerPool(p4)

	ctx, cancel := context.WithCancel(context.Background())
	config := job.NewOrchestratorConfig(250000, r)
	_, err := job.NewOrchestrator(ctx, config, c)
	if err != nil {
		fmt.Println(err)
		cancel()
	}

	for {
		current := c.CountRegisteredJobs()
		fmt.Printf("############### Jobs registered: %d\r\n", current)
		if current == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
