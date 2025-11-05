package repository

import (
	"go-task-queue-system/domain"
	"sync"
)

type MemoryRepository struct {
	tasks map[string]*domain.Task
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *MemoryRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; exists {
		return domain.ErrTaskAlreadyExists
	}

	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

func (r *MemoryRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return domain.ErrTaskNotFound
	}

	taskCopy := *task
	r.tasks[task.ID] = &taskCopy

	return nil
}

func (r *MemoryRepository) FindByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, domain.ErrTaskNotFound
	}

	taskCopy := *task
	return &taskCopy, nil
}

func (r *MemoryRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		taskCopy := *task
		tasks = append(tasks, &taskCopy)
	}

	return tasks, nil
}

func (r *MemoryRepository) FindByStatus(status domain.TaskStatus) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0)
	for _, task := range r.tasks {
		if task.Status == status {
			taskCopy := *task
			tasks = append(tasks, &taskCopy)
		}
	}

	return tasks, nil
}

func (r *MemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return domain.ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

func (r *MemoryRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tasks), nil
}

func (r *MemoryRepository) CountByStatus(status domain.TaskStatus) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, task := range r.tasks {
		if task.Status == status {
			count++
		}
	}

	return count, nil
}
