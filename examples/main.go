package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
	"github.com/corelayer/go-scheduler/pkg/task"
)

func createJob(i int) job.Job {
	schedule, _ := cron.NewSchedule("* * * * * *")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(100)
	// m := rnd.Intn(5) + 1
	tasks := []task.Task{
		task.SleepTask{
			Milliseconds: d,
		},
	}

	// if i%100 == 0 {
	// 	tasks = append(tasks, task.IntercomMessageTask{Message: fmt.Sprintf("intercom_message_%d", i)})
	// }

	d = rnd.Intn(200)
	tasks = append(tasks, task.SleepTask{Milliseconds: d})

	// if i%25 == 0 {
	return job.NewJob("Example_Job_"+fmt.Sprintf("%04d", i), schedule, 1, job.NewSequence(tasks))
	// }
	// return job.NewJob("Example_Job_"+fmt.Sprintf("%04d", i), schedule, m, job.NewSequence(tasks))
}

func handleError(err error) {
	fmt.Println(err.Error())
}

func handleMessage(msg task.IntercomMessage) {
	fmt.Println(msg.Name, msg.Content.Message)
}

func main() {
	var (
		err    error
		config job.OrchestratorConfig
	)
	c := job.NewMemoryCatalog()
	r := task.NewHandlerRepository()

	ctx, cancel := context.WithCancel(context.Background())
	config, err = job.NewOrchestratorConfig(5000, 1000, 1000, handleError, handleMessage)
	if err != nil {
		panic(err)
	}
	o := job.NewOrchestrator(c, r, config)

	err = r.RegisterHandlerPools([]*task.HandlerPool{
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

	for i < 10000 {
		i++
		if err = c.Add(createJob(i)); err != nil {
			panic(err)
		}
	}

	o.Start(ctx)

	exiting := false
	for {
		stats := o.Statistics()
		fmt.Println(stats.Job)
		if stats.Job.EnabledJobs == 0 && !exiting {
			exiting = true
			cancel()
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("--------------------------------------")
	stats := o.Statistics()

	sort.SliceStable(stats.Tasks, func(i, j int) bool {
		return stats.Tasks[i].Name < stats.Tasks[j].Name
	})

	// for _, stat := range stats.Tasks {
	// 	fmt.Println(stat)
	// } //
	// for _, v := range c.All() {
	// 	var jsonData []byte
	// 	jsonData, err = json.MarshalIndent(v, "", "\t")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(string(jsonData))
	// }

}
