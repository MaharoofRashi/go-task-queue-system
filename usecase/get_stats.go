package usecase

import "go-task-queue-system/domain"

type TaskStats struct {
	TotalTasks      int `json:"total_tasks"`
	PendingTasks    int `json:"pending_tasks"`
	ProcessingTasks int `json:"processing_tasks"`
	CompletedTasks  int `json:"completed_tasks"`
	FailedTasks     int `json:"failed_tasks"`
	CancelledTasks  int `json:"cancelled_tasks"`
	QueueSize       int `json:"queue_size"`
}

type GetStatsUseCase struct {
	repository domain.TaskRepository
	queue      TaskQueue
}

func NewGetStatsUseCase(repository domain.TaskRepository, queue TaskQueue) *GetStatsUseCase {
	return &GetStatsUseCase{
		repository: repository,
		queue:      queue,
	}
}

func (uc *GetStatsUseCase) Execute() (*TaskStats, error) {
	stats := &TaskStats{}

	total, err := uc.repository.Count()
	if err != nil {
		return nil, err
	}
	stats.TotalTasks = total

	pending, _ := uc.repository.CountByStatus(domain.TaskStatusPending)
	stats.PendingTasks = pending

	processing, _ := uc.repository.CountByStatus(domain.TaskStatusProcessing)
	stats.ProcessingTasks = processing

	completed, _ := uc.repository.CountByStatus(domain.TaskStatusCompleted)
	stats.CompletedTasks = completed

	failed, _ := uc.repository.CountByStatus(domain.TaskStatusFailed)
	stats.FailedTasks = failed

	cancelled, _ := uc.repository.CountByStatus(domain.TaskStatusCancelled)
	stats.CancelledTasks = cancelled

	stats.QueueSize = uc.queue.Size()

	return stats, nil
}
