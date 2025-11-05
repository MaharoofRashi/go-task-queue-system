package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          string                 `json:"id"`
	Type        TaskType               `json:"type"`
	Status      TaskStatus             `json:"status"`
	Priority    TaskPriority           `json:"priority"`
	Payload     map[string]interface{} `json:"payload"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string                 `json:"error,omitempty"`
	MaxRetries  int                    `json:"max_retries"`
	RetryCount  int                    `json:"retry_count"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

func NewTask(taskType TaskType, priority TaskPriority, payload map[string]interface{}) (*Task, error) {
	if !taskType.IsValid() {
		return nil, errors.New("invalid task type")
	}

	if !priority.IsValid() {
		priority = GetDefaultPriority()
	}

	if payload == nil {
		payload = make(map[string]interface{})
	}

	now := time.Now()

	return &Task{
		ID:         uuid.New().String(),
		Type:       taskType,
		Status:     TaskStatusPending,
		Priority:   priority,
		Payload:    payload,
		MaxRetries: 3,
		RetryCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (t *Task) MarkAsProcessing() {
	t.Status = TaskStatusProcessing
	now := time.Now()
	t.StartedAt = &now
	t.UpdatedAt = now
}

func (t *Task) MarkAsCompleted(result map[string]interface{}) {
	t.Status = TaskStatusCompleted
	t.Result = result
	now := time.Now()
	t.CompletedAt = &now
	t.UpdatedAt = now
}

func (t *Task) MarkAsFailed(err error) {
	t.Status = TaskStatusFailed
	if err != nil {
		t.Error = err.Error()
	}
	t.UpdatedAt = time.Now()
}

func (t *Task) MarkAsCancelled() {
	t.Status = TaskStatusCancelled
	t.UpdatedAt = time.Now()
}

func (t *Task) IncrementRetry() {
	t.RetryCount++
	t.UpdatedAt = time.Now()
}

func (t *Task) CanRetry() bool {
	return t.Status.CanRetry() && t.RetryCount < t.MaxRetries
}

func (t *Task) ShouldRetry() bool {
	return t.Status == TaskStatusFailed && t.RetryCount < t.MaxRetries
}

func (t *Task) IsInDeadLetterQueue() bool {
	return t.Status == TaskStatusFailed && t.RetryCount >= t.MaxRetries
}

func (t *Task) Validate() error {
	if t.ID == "" {
		return errors.New("task ID is required")
	}

	if !t.Type.IsValid() {
		return errors.New("invalid task type")
	}

	if !t.Status.IsValid() {
		return errors.New("invalid task status")
	}

	if !t.Priority.IsValid() {
		return errors.New("invalid task priority")
	}

	return nil
}
