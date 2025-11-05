package usecase

import "go-task-queue-system/domain"

type GetTaskUseCase struct {
	repository domain.TaskRepository
}

func NewGetTaskUseCase(repository domain.TaskRepository) *GetTaskUseCase {
	return &GetTaskUseCase{
		repository: repository,
	}
}

func (uc *GetTaskUseCase) Execute(taskID string) (*domain.Task, error) {
	if taskID == "" {
		return nil, domain.ErrTaskNotFound
	}

	task, err := uc.repository.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}
