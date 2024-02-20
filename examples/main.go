package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
	"github.com/corelayer/go-scheduler/pkg/status"
	"github.com/corelayer/go-scheduler/pkg/task"
)

func createJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("* * * * * *")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(250)
	tasks := []task.Task{
		task.TimeLogTask{},
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 1", i),
		// 	ReadInput:   false,
		// 	WriteOutput: false,
		// },
		task.EmptyTask{},
		task.SleepTask{
			Milliseconds: d,
		},
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 2", i),
		// 	ReadInput:   true,
		// 	WriteOutput: true,
		// },
		// task.SleepTask{
		// 	Milliseconds: d,
		// 	WriteOutput:  true,
		// },
		task.PrintTask{
			Message: fmt.Sprintf("Job %d", i),
		},
		task.TimeLogTask{},
	}

	return job.Job{
		Uuid:     id,
		Enabled:  true,
		Status:   status.StatusNone,
		Schedule: schedule,
		Repeat:   false,
		Name:     "Example Job " + strconv.Itoa(i),
		Tasks:    task.NewSequence(tasks),
		Intercom: task.NewIntercom(),
	}
}

func createRepeatableJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("* * * * * *")
	// rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	// d := rnd.Intn(1000)
	tasks := []task.Task{
		task.TimeLogTask{},
		task.PrintTask{
			Message: fmt.Sprintf("### Repeatable job %d - Print Task", i),
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
		Status:   status.StatusNone,
		Schedule: schedule,
		Repeat:   true,
		Name:     "### Repeatable job " + strconv.Itoa(i),
		Tasks:    task.NewSequence(tasks),
	}
}

func main() {
	c := job.NewMemoryCatalog()
	r := task.NewHandlerRepository()

	ctx, cancel := context.WithCancel(context.Background())
	config := job.NewOrchestratorConfig(2000, 2000, c, r)
	_, err := job.NewOrchestrator(ctx, config)
	if err != nil {
		fmt.Println(err)
		cancel()
	}

	err = r.RegisterHandlerPools([]*task.HandlerPool{
		task.NewHandlerPool(task.NewDefaultTimeLogTaskHandler()),
		task.NewHandlerPool(task.NewDefaultEmptyTaskHandler()),
		task.NewHandlerPool(task.NewDefaultSleepTaskHandler()),
		task.NewHandlerPool(task.NewDefaultPrintTaskHandler()),
	})
	if err != nil {
		fmt.Println(err)
	}

	i := 0
	for {
		if i < 25000 {
			fmt.Println("Adding jobs")
			if i%2500 == 0 {
				i++
				c.Register(createJob(i))
			}
			for i%2500 != 0 {
				i++
				c.Register(createJob(i))
			}
		}
		current := c.CountRegisteredJobs()
		archived := c.CountArchivedJobs()
		fmt.Printf("############### Jobs registered/archived: %d/%d\r\n", current, archived)
		if current == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
