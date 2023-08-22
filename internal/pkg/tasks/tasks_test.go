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

func (t *TestTaskRepository) GetById(id uint) (*tasks.Task, error) {
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

func (t *TestTaskRepository) DeleteById(id uint) (*tasks.Task, error) {
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
	/*
		tsks := tasks.New(tr)

		tsk, err := tsks.NewTask("First Task", time.Now().AddDate(0, 0, 15))
		if err != nil {
			t.Errorf("NewTask failed with error:%v", err)
		}

		if tsk.ID != 1 {
			t.Errorf("NewTask did not assign ID")
		}

		t.Logf("New task: %+v", tsk)

		tsk2, err := tsks.NewTask("Second Task", time.Now().AddDate(0, 0, 10))
		if err != nil {
			t.Errorf("NewTask failed with error:%v", err)
		}

		if tsk2.ID != 2 {
			t.Errorf("NewTask assigned wrong ID:%v", tsk2)
		}

		t.Logf("Second task: %+v", tsk2)

		tsklist, _ := tsks.GetAll()

		if len(tsklist) == 0 {
			t.Errorf("task list should not be empty")
		}

		tsk, err = tsks.GetById(2)
		if err != nil {
			t.Errorf("Get failed with:%v", err)
		}

		if tsk.Description != "Second Task" {
			t.Errorf("Get fetched wrong data")
		}

		tsk2, err = tsks.GetById(255)
		if err == nil {
			t.Errorf("Get retrieved invalid data:%+v", tsk2)
		}
		if !errors.Is(err, tasks.ErrNotFound) {
			t.Errorf("Expected ErrNotFound, got:%v", err)
		}

		tsk.Completed = true
		tsk2, err = tsks.Update(tsk)
		if err != nil {
			t.Errorf("Update failed with:%v", err)
		}
		if tsk.Description != tsk2.Description || tsk.Deadline != tsk2.Deadline {
			t.Errorf("wrong data got updated")
		}

		if !tsk2.Completed {
			t.Errorf("data was not updated")
		}

		tsk, err = tsks.GetById(2)
		if err != nil || !tsk.Completed {
			t.Errorf("data was not updated")
		}

		_, err = tsks.DeleteById(255)
		if err == nil || !errors.Is(err, tasks.ErrNotFound) {
			t.Errorf("unexpected error behaviour:%v", err)
		}

		tsk, err = tsks.DeleteById(2)
		if err != nil {
			t.Errorf("DeleteById failed with:%v", err)
		}

		t.Logf("Deleted data:%+v", tsk)

		tsk2, err = tsks.GetById(2)
		if !errors.Is(err, tasks.ErrNotFound) {
			t.Errorf("delete did not work. Data found:%+v", tsk2)
		}
	*/
}
