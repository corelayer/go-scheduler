package main

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/google/uuid"

	"github.com/corelayer/go-scheduler/pkg/cron"
	"github.com/corelayer/go-scheduler/pkg/job"
)

type TimeLogTask struct {
	timestamp time.Time
}

func (t TimeLogTask) WriteToPipeline() bool {
	return true
}

type TimeLogTaskHandler struct{}

func (h TimeLogTaskHandler) Execute(t job.Task, pipeline chan interface{}) job.Task {
	task := t.(TimeLogTask)
	task.timestamp = time.Now()

	select {
	case data := <-pipeline:
		if task.WriteToPipeline() {
			pipeline <- data
		}
	default:
	}

	return task
}

func (h TimeLogTaskHandler) GetTaskType() reflect.Type {
	return reflect.TypeOf(TimeLogTask{})
}

func createJob(i int) job.Job {
	id, _ := uuid.NewUUID()
	schedule, _ := cron.NewSchedule("@everysecond")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := rnd.Intn(1000)
	tasks := []job.Task{
		TimeLogTask{},
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
		TimeLogTask{},
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

	p1 := job.NewTaskHandlerPool(job.PrintTaskHandler{}, 10)
	p2 := job.NewTaskHandlerPool(job.SleepTaskHandler{}, 5)
	p3 := job.NewTaskHandlerPool(TimeLogTaskHandler{}, 100)

	r := job.NewTaskHandlerRepository()
	r.RegisterTaskHandlerPool(p1)
	r.RegisterTaskHandlerPool(p2)
	r.RegisterTaskHandlerPool(p3)

	timeout, _ := time.ParseDuration("20s")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	config := job.NewOrchestratorConfig(10, r)
	_, err := job.NewOrchestrator(ctx, config, c)
	if err != nil {
		fmt.Println(err)
		cancel()
	}

	// time.Sleep(30 * time.Second)
	time.Sleep(5 * time.Second)

	jobs := c.GetAllJobs()
	for i, j := range jobs {
		if j.Status == job.StatusCompleted {
			startTime := j.Tasks.Tasks[0].(TimeLogTask).timestamp
			endTime := j.Tasks.Tasks[6].(TimeLogTask).timestamp
			runTime := endTime.Sub(startTime)
			fmt.Println(i, j.Status.String(), runTime.String())
		}
	}
}
