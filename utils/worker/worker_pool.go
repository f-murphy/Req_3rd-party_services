package utils

import (
	"req3rdPartyServices/models"
	utils "req3rdPartyServices/utils/executor"
	"sync"
)

type Worker struct {
	taskChan chan *models.Task
	quitChan chan bool
	wg       sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{
		taskChan: make(chan *models.Task),
		quitChan: make(chan bool),
	}
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			select {
			case task := <-w.taskChan:
				utils.ExecuteTask(task)
			case <-w.quitChan:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quitChan)
	w.wg.Wait()
}

func (w *Worker) AddTask(task *models.Task) {
	w.taskChan <- task
}
