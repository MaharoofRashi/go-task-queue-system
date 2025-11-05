package http

import "go-task-queue-system/domain"

type SubmitTaskRequest struct {
	Type     string                 `json:"type"`
	Priority string                 `json:"priority,omitempty"`
	Payload  map[string]interface{} `json:"payload"`
}

type TaskResponse struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Status      string                 `json:"status"`
	Priority    string                 `json:"priority"`
	Payload     map[string]interface{} `json:"payload"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	MaxRetries  int                    `json:"max_retries"`
	RetryCount  int                    `json:"retry_count"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	StartedAt   *string                `json:"started_at,omitempty"`
	CompletedAt *string                `json:"completed_at,omitempty"`
}

type TaskListResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
	Total int             `json:"total"`
}

type StatsResponse struct {
	TotalTasks      int `json:"total_tasks"`
	PendingTasks    int `json:"pending_tasks"`
	ProcessingTasks int `json:"processing_tasks"`
	CompletedTasks  int `json:"completed_tasks"`
	FailedTasks     int `json:"failed_tasks"`
	CancelledTasks  int `json:"cancelled_tasks"`
	QueueSize       int `json:"queue_size"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

type WorkerStatusResponse struct {
	WorkerCount int    `json:"worker_count"`
	Timeout     string `json:"timeout"`
}

func ToTaskResponse(task *domain.Task) *TaskResponse {
	response := &TaskResponse{
		ID:         task.ID,
		Type:       task.Type.String(),
		Status:     task.Status.String(),
		Priority:   task.Priority.String(),
		Payload:    task.Payload,
		Result:     task.Result,
		Error:      task.Error,
		MaxRetries: task.MaxRetries,
		RetryCount: task.RetryCount,
		CreatedAt:  task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if task.StartedAt != nil {
		startedAt := task.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		response.StartedAt = &startedAt
	}

	if task.CompletedAt != nil {
		completedAt := task.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	return response
}

func ToTaskListResponse(tasks []*domain.Task) *TaskListResponse {
	taskResponses := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = ToTaskResponse(task)
	}

	return &TaskListResponse{
		Tasks: taskResponses,
		Total: len(tasks),
	}
}
