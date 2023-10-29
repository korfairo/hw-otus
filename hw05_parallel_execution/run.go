package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var numErr int64
	errLimit := int64(m)

	taskCh := make(chan Task)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			worker(taskCh, &numErr)
			wg.Done()
		}()
	}

	var err error
	for _, task := range tasks {
		if atomic.LoadInt64(&numErr) >= errLimit {
			err = ErrErrorsLimitExceeded
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()
	return err
}

func worker(taskCh chan Task, numErr *int64) {
	for task := range taskCh {
		if task == nil {
			continue
		}

		if err := task(); err != nil {
			atomic.AddInt64(numErr, 1)
		}
	}
}
