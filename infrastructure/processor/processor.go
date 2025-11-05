package processor

import (
	"context"
	"go-task-queue-system/domain"
)

type TaskProcessor interface {
	Process(ctx context.Context, task *domain.Task) (map[string]interface{}, error)

	CanProcess(taskType domain.TaskType) bool
}

type ProcessorRegistry struct {
	processors map[domain.TaskType]TaskProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{
		processors: make(map[domain.TaskType]TaskProcessor),
	}
}

func (r *ProcessorRegistry) Register(taskType domain.TaskType, processor TaskProcessor) {
	r.processors[taskType] = processor
}

func (r *ProcessorRegistry) GetProcessor(taskType domain.TaskType) (TaskProcessor, bool) {
	processor, exists := r.processors[taskType]
	return processor, exists
}

func (r *ProcessorRegistry) HasProcessor(taskType domain.TaskType) bool {
	_, exists := r.processors[taskType]
	return exists
}
