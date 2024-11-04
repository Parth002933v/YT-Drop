package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Task interface {
	Process(context.Context)
}

type Pool struct {
	Tasks         []Task
	MaxConcurrent int
	taskChan      chan Task
	wg            sync.WaitGroup
}

func (wp *Pool) worker(ctx context.Context) {

	for task := range wp.taskChan {
		select {
		case <-ctx.Done():
			fmt.Println(`worker received canal signal, stopping... `)
			return
		default:
			task.Process(ctx)
			wp.wg.Done()
		}
	}
}
func (wp *Pool) Run() {
	wp.taskChan = make(chan Task, len(wp.Tasks))
	wp.wg.Add(len(wp.Tasks))

	ctx, cancal := context.WithCancel(context.Background())
	defer cancal()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		fmt.Println("Interrupt signal received,canceling tasks... ")
		cancal()
		close(wp.taskChan)
		os.Exit(1)
	}()

	for i := 0; i < wp.MaxConcurrent; i++ {
		go wp.worker(ctx)
	}

	for _, task := range wp.Tasks {
		select {
		case <-ctx.Done():
			fmt.Println("context canceled before all tasks are queuedj")
			break
		default:
			wp.taskChan <- task
		}
	}

	close(wp.taskChan)
	wp.wg.Wait()
}
