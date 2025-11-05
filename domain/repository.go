package domain

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
)

type TaskRepository interface {
	Save(task *Task) error

	Update(task *Task) error

	FindByID(id string) (*Task, error)

	FindAll() ([]*Task, error)

	FindByStatus(status TaskStatus) ([]*Task, error)

	Delete(id string) error

	Count() (int, error)

	CountByStatus(status TaskStatus) (int, error)
}
