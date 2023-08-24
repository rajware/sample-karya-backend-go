package tasks_test

import (
	"testing"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks/taskstest"
	"golang.org/x/exp/slices"
)

type TestTaskRepository struct {
	nextID uint
	tasks  []*tasks.Task
}

func (t *TestTaskRepository) GetAll() ([]*tasks.Task, error) {
	return t.tasks, nil
}

func (t *TestTaskRepository) GetByID(id uint) (*tasks.Task, error) {
	index := slices.IndexFunc(t.tasks, func(t *tasks.Task) bool {
		if t != nil {
			return t.ID == id
		}

		return false
	})
	if index == -1 {
		return nil, tasks.ErrNotFound
	}
	return t.tasks[index], nil
}

func (t *TestTaskRepository) Add(nt *tasks.Task) (*tasks.Task, error) {
	nt.ID = t.nextID
	t.tasks = append(t.tasks, nt)
	t.nextID++
	return nt, nil
}

func (t *TestTaskRepository) Update(ut *tasks.Task) (*tasks.Task, error) {
	if ut == nil {
		return nil, tasks.ErrNotFound
	}

	index := slices.IndexFunc(t.tasks, func(t *tasks.Task) bool {
		return t.ID == ut.ID
	})
	if index == -1 {
		return nil, tasks.ErrNotFound
	}
	t.tasks[index].Description = ut.Description
	t.tasks[index].Deadline = ut.Deadline
	t.tasks[index].Completed = ut.Completed

	return t.tasks[index], nil
}

func (t *TestTaskRepository) DeleteByID(id uint) (*tasks.Task, error) {
	index := slices.IndexFunc(t.tasks, func(t *tasks.Task) bool {
		return t.ID == id
	})
	if index == -1 {
		return nil, tasks.ErrNotFound
	}

	deletedtask := *t.tasks[index]
	t.tasks[index] = nil
	t.tasks = slices.Delete(t.tasks, index, index)

	return &deletedtask, nil
}

func TestTasks(t *testing.T) {
	tr := &TestTaskRepository{
		nextID: 1,
		tasks:  []*tasks.Task{},
	}

	taskstest.TestTaskRepository(t, tr)
}
