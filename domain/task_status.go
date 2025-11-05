package domain

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusPending, TaskStatusProcessing, TaskStatusCompleted, TaskStatusFailed, TaskStatusCancelled:
		return true
	default:
		return false
	}
}

func (s TaskStatus) String() string {
	return string(s)
}

func (s TaskStatus) IsFinal() bool {
	return s == TaskStatusCompleted || s == TaskStatusCancelled
}

func (s TaskStatus) CanRetry() bool {
	return s == TaskStatusFailed
}
