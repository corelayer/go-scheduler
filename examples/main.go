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
	m := rnd.Intn(5) + 1
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

	if i%100 == 0 {
		tasks = append(tasks, task.IntercomMessageTask{Message: fmt.Sprintf("intercom_message_%d", i)})
	}

	return job.NewJob("Example_Job_"+strconv.Itoa(i), schedule, m, task.NewSequence(tasks))
}

func handleError(err error) {
	fmt.Println(err.Error())
}

func handleMessage(msg task.IntercomMessage) {
	fmt.Println(msg.Name, msg.Content.Message)
}

func main() {
	c := job.NewMemoryCatalog()
	r := task.NewHandlerRepository()

	ctx, cancel := context.WithCancel(context.Background())
	config := job.NewOrchestratorConfig(10, 1, handleError, nil)
	o := job.NewOrchestrator(c, r, config)

	err := r.RegisterHandlerPools([]*task.HandlerPool{
		task.NewHandlerPool(task.NewDefaultTimeLogTaskHandler()),
		task.NewHandlerPool(task.NewDefaultEmptyTaskHandler()),
		task.NewHandlerPool(task.NewDefaultSleepTaskHandler()),
		task.NewHandlerPool(task.NewDefaultPrintTaskHandler()),
		task.NewHandlerPool(task.NewDefaultIntercomMessageTaskHandler()),
	})
	if err != nil {
		fmt.Println(err)
	}

	i := 0

	for i < 3000 {
		i++
		if err = c.Add(createJob(i)); err != nil {
			panic(err)
		}
	}

	go o.Start(ctx)

	for c.HasEnabledJobs() {
		time.Sleep(250 * time.Millisecond)
	}
	cancel()
	i = 0
	for _, j := range c.All() {
		if j.IsEnabled() {
			fmt.Println(j.Name)
		} else {
			results := j.Results()
			fmt.Println(j.Name)
			for _, r := range results {
				fmt.Println("\t", r.RunTime, r.Status)
			}
		}
		i++
	}
	fmt.Println("jobs retrieved", i)

}
