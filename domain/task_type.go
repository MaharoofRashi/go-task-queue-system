package domain

type TaskType string

const (
	TaskTypeEmail            TaskType = "email"
	TaskTypeImageProcessing  TaskType = "image_processing"
	TaskTypeReportGeneration TaskType = "report_generation"
)

func (t TaskType) IsValid() bool {
	switch t {
	case TaskTypeEmail, TaskTypeImageProcessing, TaskTypeReportGeneration:
		return true
	default:
		return false
	}
}

func (t TaskType) String() string {
	return string(t)
}
