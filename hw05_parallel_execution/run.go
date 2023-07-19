package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		wg       sync.WaitGroup
		errCount int64
		taskCh   = make(chan Task)
		errCh    = make(chan struct{}, len(tasks))
	)

	defer func() {
		close(taskCh)
		wg.Wait()
		close(errCh)
	}()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				t, ok := <-taskCh
				if !ok {
					return
				}
				if err := t(); err != nil {
					errCh <- struct{}{}
				}
			}
		}()
	}

	for _, t := range tasks {
		if m > 0 && (errCount == int64(m)) {
			return ErrErrorsLimitExceeded
		}
		taskCh <- t
		select {
		case <-errCh:
			errCount++
		default:
		}
	}

	return nil
}
