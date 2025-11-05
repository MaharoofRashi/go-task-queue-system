package worker

import (
	"context"
	"fmt"
	"go-task-queue-system/domain"
	"go-task-queue-system/infrastructure/processor"
	"log"
	"time"
)

type Worker struct {
	id                int
	taskQueue         <-chan *domain.Task
	repository        domain.TaskRepository
	processorRegistry *processor.ProcessorRegistry
	quit              chan bool
	timeout           time.Duration
}

func NewWorker(
	id int,
	taskQueue <-chan *domain.Task,
	repository domain.TaskRepository,
	processorRegistry *processor.ProcessorRegistry,
	timeout time.Duration,
) *Worker {
	return &Worker{
		id:                id,
		taskQueue:         taskQueue,
		repository:        repository,
		processorRegistry: processorRegistry,
		quit:              make(chan bool),
		timeout:           timeout,
	}
}

func (w *Worker) Start() {
	log.Printf("🚀 Worker %d started", w.id)

	for {
		select {
		case task, ok := <-w.taskQueue:
			if !ok {
				log.Printf("⛔ Worker %d: task queue closed", w.id)
				return
			}
			w.processTask(task)

		case <-w.quit:
			log.Printf("⛔ Worker %d: received quit signal", w.id)
			return
		}
	}
}

func (w *Worker) Stop() {
	log.Printf("🛑 Stopping worker %d", w.id)
	w.quit <- true
}

func (w *Worker) processTask(task *domain.Task) {
	log.Printf("⚙️  Worker %d: picked up task %s (type: %s)", w.id, task.ID, task.Type)

	task.MarkAsProcessing()
	if err := w.repository.Update(task); err != nil {
		log.Printf("❌ Worker %d: failed to update task status: %v", w.id, err)
		return
	}

	proc, exists := w.processorRegistry.GetProcessor(task.Type)
	if !exists {
		log.Printf("❌ Worker %d: no processor found for task type %s", w.id, task.Type)
		task.MarkAsFailed(fmt.Errorf("no processor found for task type: %s", task.Type))
		w.repository.Update(task)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), w.timeout)
	defer cancel()

	result, err := proc.Process(ctx, task)

	if err != nil {
		log.Printf("❌ Worker %d: task %s failed: %v", w.id, task.ID, err)
		task.MarkAsFailed(err)
		task.IncrementRetry()

		if task.ShouldRetry() {
			log.Printf("🔄 Worker %d: task %s will be retried (attempt %d/%d)",
				w.id, task.ID, task.RetryCount, task.MaxRetries)
			// In a real system, have to re-enqueue the task here
			// For now, just mark it as failed and it stays failed
		} else if task.IsInDeadLetterQueue() {
			log.Printf("☠️  Worker %d: task %s moved to dead letter queue (max retries exceeded)",
				w.id, task.ID)
		}

		w.repository.Update(task)
		return
	}

	log.Printf("✅ Worker %d: task %s completed successfully", w.id, task.ID)
	task.MarkAsCompleted(result)
	w.repository.Update(task)
}
