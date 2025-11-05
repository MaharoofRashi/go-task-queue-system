package usecase

import "go-task-queue-system/domain"

type ListTasksUseCase struct {
	repository domain.TaskRepository
}

func NewListTasksUseCase(repository domain.TaskRepository) *ListTasksUseCase {
	return &ListTasksUseCase{
		repository: repository,
	}
}

func (uc *ListTasksUseCase) Execute() ([]*domain.Task, error) {
	tasks, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (uc *ListTasksUseCase) ExecuteByStatus(status domain.TaskStatus) ([]*domain.Task, error) {
	if !status.IsValid() {
		return nil, domain.ErrInvalidTaskStatus
	}

	tasks, err := uc.repository.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
