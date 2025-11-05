package domain

import "errors"

var (
	ErrInvalidTaskType = errors.New("invalid task type")

	ErrInvalidTaskStatus = errors.New("invalid task status")

	ErrInvalidTaskPriority = errors.New("invalid task priority")

	ErrTaskCannotBeRetried = errors.New("task cannot be retried")

	ErrTaskAlreadyProcessing = errors.New("task is already being processed")

	ErrTaskAlreadyCompleted = errors.New("task is already completed")

	ErrEmptyPayload = errors.New("task payload cannot be empty")
)
