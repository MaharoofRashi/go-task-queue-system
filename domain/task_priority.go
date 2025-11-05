package domain

type TaskPriority string

const (
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityLow    TaskPriority = "low"
)

func (p TaskPriority) IsValid() bool {
	switch p {
	case TaskPriorityHigh, TaskPriorityMedium, TaskPriorityLow:
		return true
	default:
		return false
	}
}

func (p TaskPriority) String() string {
	return string(p)
}

func GetDefaultPriority() TaskPriority {
	return TaskPriorityMedium
}
