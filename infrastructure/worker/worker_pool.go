package worker

import (
	"go-task-queue-system/domain"
	"go-task-queue-system/infrastructure/processor"
	"log"
	"sync"
	"time"
)

type WorkerPool struct {
	workers           []*Worker
	workerCount       int
	taskQueue         <-chan *domain.Task
	repository        domain.TaskRepository
	processorRegistry *processor.ProcessorRegistry
	timeout           time.Duration
	wg                sync.WaitGroup
}

func NewWorkerPool(
	workerCount int,
	taskQueue <-chan *domain.Task,
	repository domain.TaskRepository,
	processorRegistry *processor.ProcessorRegistry,
	timeout time.Duration,
) *WorkerPool {
	return &WorkerPool{
		workers:           make([]*Worker, 0, workerCount),
		workerCount:       workerCount,
		taskQueue:         taskQueue,
		repository:        repository,
		processorRegistry: processorRegistry,
		timeout:           timeout,
	}
}

func (wp *WorkerPool) Start() {
	log.Printf("🚀 Starting worker pool with %d workers", wp.workerCount)

	for i := 1; i <= wp.workerCount; i++ {
		worker := NewWorker(
			i,
			wp.taskQueue,
			wp.repository,
			wp.processorRegistry,
			wp.timeout,
		)

		wp.workers = append(wp.workers, worker)

		wp.wg.Add(1)
		go func(w *Worker) {
			defer wp.wg.Done()
			w.Start()
		}(worker)
	}

	log.Printf("✅ Worker pool started successfully")
}

func (wp *WorkerPool) Stop() {
	log.Printf("🛑 Stopping worker pool...")

	for _, worker := range wp.workers {
		worker.Stop()
	}

	wp.wg.Wait()

	log.Printf("✅ Worker pool stopped")
}

func (wp *WorkerPool) GetWorkerCount() int {
	return wp.workerCount
}

func (wp *WorkerPool) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"worker_count": wp.workerCount,
		"timeout":      wp.timeout.String(),
	}
}
