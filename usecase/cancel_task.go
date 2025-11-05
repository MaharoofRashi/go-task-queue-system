package usecase

import (
	"errors"
	"go-task-queue-system/domain"
)

type CancelTaskUseCase struct {
	repository domain.TaskRepository
}

func NewCancelTaskUseCase(repository domain.TaskRepository) *CancelTaskUseCase {
	return &CancelTaskUseCase{
		repository: repository,
	}
}

func (uc *CancelTaskUseCase) Execute(taskID string) error {
	if taskID == "" {
		return domain.ErrTaskNotFound
	}

	task, err := uc.repository.FindByID(taskID)
	if err != nil {
		return err
	}

	if task.Status != domain.TaskStatusPending {
		return errors.New("only pending tasks can be cancelled")
	}

	task.MarkAsCancelled()

	if err := uc.repository.Update(task); err != nil {
		return err
	}

	return nil
}
