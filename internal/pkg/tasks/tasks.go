package tasks

import (
	"errors"
	"time"
)

var ErrNoRepo = errors.New("data store not initialized")
var ErrNotFound = errors.New("data not found")
var ErrNotUpdated = errors.New("update failed")

type Task struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed"`
}

type TaskRepository interface {
	GetAll() ([]*Task, error)
	GetById(id uint) (*Task, error)
	Add(*Task) (*Task, error)
	Update(*Task) (*Task, error)
	DeleteById(id uint) (*Task, error)
}

type Tasks struct {
	repo TaskRepository
}

func (t *Tasks) validaterepo() bool {
	return t.repo != nil
}

func (t *Tasks) NewTask(description string, deadline time.Time) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.Add(&Task{
		Description: description,
		Deadline:    deadline,
		Completed:   false,
	})
}

func (t *Tasks) GetAll() ([]*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.GetAll()
}

func (t *Tasks) GetById(id uint) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.GetById(id)
}

func (t *Tasks) Update(tsk *Task) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.Update(tsk)
}

func (t *Tasks) DeleteById(id uint) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.DeleteById(id)
}

func New(r TaskRepository) *Tasks {
	return &Tasks{
		repo: r,
	}
}
