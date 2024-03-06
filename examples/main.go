package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
	"github.com/corelayer/go-scheduler/pkg/task"
)

func createJob(i int) job.Job {
	schedule, _ := cron.NewSchedule("* * * * * *")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(100)
	tasks := []task.Task{
		// task.TimeLogTask{},
		// task.PrintTask{
		// 	Message:     fmt.Sprintf("Job %d - Task 1", i),
		// 	ReadInput:   false,
		// 	WriteOutput: false,
		// },
		// task.EmptyTask{},
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
		// task.PrintTask{
		// 	Message: fmt.Sprintf("Job %d", i),
		// },
		// task.TimeLogTask{},
	}

	return job.NewJob("Example Job "+strconv.Itoa(i), schedule, 2, task.NewSequence(tasks))
}

func handleError(err error) {
	fmt.Println(err.Error())
}

func handleMessage(msg task.Message) {
	fmt.Println("Message", msg.Task, msg.Data)
}

func main() {
	c := job.NewMemoryCatalog()
	r := task.NewHandlerRepository()

	ctx, cancel := context.WithCancel(context.Background())
	config := job.NewOrchestratorConfig(20, 1, handleError, handleMessage)
	o := job.NewOrchestrator(c, r, config)

	err := r.RegisterHandlerPools([]*task.HandlerPool{
		task.NewHandlerPool(task.NewDefaultTimeLogTaskHandler()),
		task.NewHandlerPool(task.NewDefaultEmptyTaskHandler()),
		task.NewHandlerPool(task.NewDefaultSleepTaskHandler()),
		task.NewHandlerPool(task.NewDefaultPrintTaskHandler()),
	})
	if err != nil {
		fmt.Println(err)
	}

	i := 0
	for i < 25000 {
		if i%2500 == 0 {
			fmt.Println("Adding jobs", i)
			i++
			if err = c.Add(createJob(i)); err != nil {
				panic(err)
			}
			// time.Sleep(1 * time.Second)
		}
		for i%2500 != 0 {
			i++
			if err = c.Add(createJob(i)); err != nil {
				panic(err)
			}
		}
	}
	// current := c.CountRegisteredJobs()
	// archived := c.CountArchivedJobs()
	// fmt.Printf("############### Jobs registered/archived: %d/%d\r\n", current, archived)
	// if current == 0 {
	// 	break
	// }
	// time.Sleep(1 * time.Second)

	go o.Start(ctx)

	for c.HasEnabledJobs() {
		fmt.Println("waiting")
		time.Sleep(250 * time.Millisecond)
	}
	cancel()

	for _, j := range c.All() {
		if j.IsEnabled() {
			fmt.Println(j)
		} else {
			fmt.Println(j.Results())
		}
	}
}
