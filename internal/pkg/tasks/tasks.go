package tasks

import (
	"errors"
	"time"
)

// Well-known errors that can occur during Tasks manipulation
var (
	ErrNoRepo     = errors.New("data store not initialized")
	ErrNotFound   = errors.New("data not found")
	ErrNotUpdated = errors.New("update failed")
)

// Task represents a task.
type Task struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed"`
}

// TaskRepository represents any data store that can store Tasks.
type TaskRepository interface {
	GetAll() ([]*Task, error)          // Get all Tasks, undefined order
	GetByID(id uint) (*Task, error)    // Get Task by ID, or error
	Add(*Task) (*Task, error)          // Add Task, providing auto generated ID
	Update(*Task) (*Task, error)       // Update Task. Task should exist
	DeleteByID(id uint) (*Task, error) // Delete Task by ID, or error. Task should exist
}

// Tasks represents a collection of tasks stored and retrieved from
// a repository.
type Tasks struct {
	repo TaskRepository
}

func (t *Tasks) validaterepo() bool {
	return t.repo != nil
}

// NewTask creates a new Task with an auto-generated ID, the given
// Description and Deadline, and Completed status of false, and
// stores it in the repository.
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

// GetAll returns all Tasks stored in the repository.
func (t *Tasks) GetAll() ([]*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.GetAll()
}

// GetByID returns the Task with the given ID, or an error.
func (t *Tasks) GetByID(id uint) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.GetByID(id)
}

// Update updates the given task in the repository, and
// returns the updated task or an error. The task should
// exist in the repository.
func (t *Tasks) Update(tsk *Task) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.Update(tsk)
}

// DeleteByID deletes the Task with the given ID from the
// repository, and returns the deleted task or an error.
// The Task should exist in the repository.
func (t *Tasks) DeleteByID(id uint) (*Task, error) {
	if !t.validaterepo() {
		return nil, ErrNoRepo
	}

	return t.repo.DeleteByID(id)
}

// New returns a new Tasks domain context object. It needs
// to be provided with an implementation of TaskRepository
// , which will store and retrieve Tasks in some storage
// backend.
func New(r TaskRepository) *Tasks {
	return &Tasks{
		repo: r,
	}
}
