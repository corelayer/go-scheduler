package main

import (
	"context"
	"encoding/json"
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
	m := rnd.Intn(50) + 1
	tasks := []task.Task{
		task.SleepTask{
			Milliseconds: d,
		},
	}

	if i%2 == 0 {
		tasks = append(tasks, task.IntercomMessageTask{Message: fmt.Sprintf("intercom_message_%d", i)})
	}

	d = rnd.Intn(200)
	tasks = append(tasks, task.SleepTask{Milliseconds: d})

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
	config := job.NewOrchestratorConfig(2000, 100, handleError, handleMessage)
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

	for i < 4 {
		i++
		if err = c.Add(createJob(i)); err != nil {
			panic(err)
		}
	}

	o.Start(ctx)

	exiting := false
	for {
		if !o.IsStarted() {
			break
		}
		stats := o.Statistics()
		fmt.Println("Active", stats.ActiveJobs, "-", "Enabled", stats.EnabledJobs)

		if stats.EnabledJobs == 0 && !exiting {
			exiting = true
			cancel()
		}
		time.Sleep(1 * time.Second)
	}

	for _, v := range c.All() {
		var jsonData []byte
		jsonData, err = json.MarshalIndent(v, "", "\t")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(jsonData))
	}

}
