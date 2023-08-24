package taskstest

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
)

func dataDirectoryPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	datadir := filepath.Join(wd, "data")
	return datadir, nil
}

// SetupDataDirectory creates a subdirectory called 'data' under the
// current working directory and returns the absolute path, or an error.
func SetupDataDirectory() (string, error) {
	datadir, err := dataDirectoryPath()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(datadir, 0755)
	if err != nil {
		return "", err
	}

	return datadir, nil
}

// RemoveDataDirectory deletes a subdirectory called 'data' under the
// current working directory, or returns an error.
func RemoveDataDirectory() error {
	datadir, err := dataDirectoryPath()
	if err != nil {
		return err
	}
	err = os.RemoveAll(datadir)
	if err != nil {
		return err
	}

	return nil
}

// TestTaskRepository runs a full suite of tests on a tasks.TestTaskRepository
// implementation.
func TestTaskRepository(t *testing.T, tr tasks.TaskRepository) {
	tsks := tasks.New(tr)

	tsk1, err := tsks.NewTask("First Task", time.Now().AddDate(0, 0, 15))
	if err != nil {
		t.Errorf("NewTask failed with error:%v", err)
	}

	tsk1ID := tsk1.ID

	if tsk1ID == 0 {
		t.Errorf("NewTask did not assign ID")
	}

	t.Logf("New task: %+v", tsk1)

	tsk2, err := tsks.NewTask("Second Task", time.Now().AddDate(0, 0, 10))
	if err != nil {
		t.Errorf("NewTask failed with error:%v", err)
	}

	tsk2ID := tsk2.ID

	if tsk2ID == 0 || tsk2ID == tsk1ID {
		t.Errorf("NewTask assigned wrong ID:%v", tsk2)
	}

	t.Logf("Second task: %+v", tsk2)

	tsklist, _ := tsks.GetAll()

	if len(tsklist) == 0 {
		t.Errorf("task list should not be empty")
	}

	tsk1, err = tsks.GetByID(tsk2ID)
	if err != nil {
		t.Errorf("Get failed with:%v", err)
	}

	if tsk1.Description != "Second Task" {
		t.Errorf("Get fetched wrong data")
	}

	tsk2, err = tsks.GetByID(255)
	if err == nil {
		t.Errorf("Get retrieved invalid data:%+v", tsk2)
	}
	if !errors.Is(err, tasks.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got:%v", err)
	}

	tsk1.Completed = true
	tsk2, err = tsks.Update(tsk1)
	if err != nil {
		t.Errorf("Update failed with:%v", err)
	}
	if tsk1.Description != tsk2.Description || tsk1.Deadline != tsk2.Deadline {
		t.Errorf("wrong data got updated")
	}

	if !tsk2.Completed {
		t.Errorf("data was not updated")
	}

	tsk1, err = tsks.GetByID(tsk2ID)
	if err != nil || !tsk1.Completed {
		t.Errorf("data was not updated")
	}

	_, err = tsks.DeleteByID(255)
	if err == nil || !errors.Is(err, tasks.ErrNotFound) {
		t.Errorf("unexpected error behaviour:%v", err)
	}

	tsk1, err = tsks.DeleteByID(tsk2ID)
	if err != nil {
		t.Errorf("DeleteById failed with:%v", err)
	}

	t.Logf("Deleted data:%+v", tsk1)

	tsk2, err = tsks.GetByID(tsk2ID)
	if !errors.Is(err, tasks.ErrNotFound) {
		t.Errorf("delete did not work. Data found:%+v", tsk2)
	}
}
