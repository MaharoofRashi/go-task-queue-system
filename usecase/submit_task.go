package usecase

import (
	"errors"
	"go-task-queue-system/domain"
)

type SubmitTaskUseCase struct {
	repository domain.TaskRepository
	queue      TaskQueue
}

type TaskQueue interface {
	Enqueue(task *domain.Task) error
	Size() int
}

func NewSubmitTaskUseCase(repository domain.TaskRepository, queue TaskQueue) *SubmitTaskUseCase {
	return &SubmitTaskUseCase{
		repository: repository,
		queue:      queue,
	}
}

func (uc *SubmitTaskUseCase) Execute(taskType domain.TaskType, priority domain.TaskPriority, payload map[string]interface{}) (*domain.Task, error) {
	if !taskType.IsValid() {
		return nil, domain.ErrInvalidTaskType
	}

	if payload == nil || len(payload) == 0 {
		return nil, domain.ErrEmptyPayload
	}

	task, err := domain.NewTask(taskType, priority, payload)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Save(task); err != nil {
		return nil, err
	}

	if err := uc.queue.Enqueue(task); err != nil {
		// If enqueue fails, later probably mark the task as failed
		// For now, just return the error
		return nil, errors.New("failed to enqueue task: " + err.Error())
	}

	return task, nil
}
