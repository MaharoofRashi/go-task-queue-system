package queue

import (
	"errors"
	"go-task-queue-system/domain"
)

type ChannelQueue struct {
	taskChan chan *domain.Task
	capacity int
}

func NewChannelQueue(capacity int) *ChannelQueue {
	return &ChannelQueue{
		taskChan: make(chan *domain.Task, capacity),
		capacity: capacity,
	}
}

func (q *ChannelQueue) Enqueue(task *domain.Task) error {
	select {
	case q.taskChan <- task:
		return nil
	default:
		return errors.New("queue is full")
	}
}

func (q *ChannelQueue) Dequeue() (*domain.Task, error) {
	task, ok := <-q.taskChan
	if !ok {
		return nil, errors.New("queue is closed")
	}
	return task, nil
}

func (q *ChannelQueue) Size() int {
	return len(q.taskChan)
}

func (q *ChannelQueue) GetChannel() <-chan *domain.Task {
	return q.taskChan
}

func (q *ChannelQueue) Close() {
	close(q.taskChan)
}

func (q *ChannelQueue) Capacity() int {
	return q.capacity
}

func (q *ChannelQueue) IsFull() bool {
	return len(q.taskChan) >= q.capacity
}

func (q *ChannelQueue) IsEmpty() bool {
	return len(q.taskChan) == 0
}
