package workerpool

import (
	"sync"
)

type WorkerPool struct {
	workers chan func()
	queue   chan func()
	wg      sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{
		workers: make(chan func(), size),
		queue:   make(chan func()),
	}
	wp.start()
	wp.startWorkers()
	return wp
}

func (wp *WorkerPool) start() {
	go func() {
		for task := range wp.queue {
			wp.workers <- task
		}
	}()
}

func (wp *WorkerPool) startWorkers() {
	for i := 0; i < cap(wp.workers); i++ {
		wp.wg.Add(1)
		go func() {
			for task := range wp.workers {
				task()
			}
			wp.wg.Done()
		}()
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.queue <- task
}

func (wp *WorkerPool) Close() {
	close(wp.queue)
	wp.wg.Wait()
}
