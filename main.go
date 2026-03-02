package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Task interface {
	Execute(ctx context.Context) error
}

type EmailTask struct {
	ID int
}

func (t *EmailTask) Execute(ctx context.Context) error {
	fmt.Println("executing task", t.ID)
	return nil
}

func worker(id int, tasks <-chan Task, results chan<- int, r *rate.Limiter) {
	for task := range tasks {
		r.Wait(context.Background())
		task.Execute(context.Background())
		fmt.Println("worker", id, "started job", task)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", task)
		results <- id
	}
}

func main() {

	var wg sync.WaitGroup
	var r = rate.NewLimiter(rate.Every(time.Second), 1)

	const numTasks = 5
	tasks := make(chan Task, numTasks)
	results := make(chan int, numTasks)

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, tasks, results, r)
		}(w)
	}
	for j := 1; j <= numTasks; j++ {
		tasks <- &EmailTask{ID: j}
	}
	close(tasks)

	wg.Wait()
	close(results)
	for a := 1; a <= numTasks; a++ {
		fmt.Println("result: ", <-results)
	}
}
