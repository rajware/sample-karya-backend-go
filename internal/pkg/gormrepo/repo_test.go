package gormrepo_test

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks/taskstest"
	"gorm.io/gorm"
)

func TestGORMTaskRepository(t *testing.T) {
	datadir, err := taskstest.SetupDataDirectory()
	if err != nil {
		t.Logf("data directory creation failed with:%v", err)
		t.FailNow()
	}
	defer taskstest.RemoveDataDirectory()

	datafile := filepath.Join(datadir, "test.db")

	t.Logf("Opening data file %v", datafile)
	d := sqlite.Open(datafile)
	tr := gormrepo.New(d, &gorm.Config{})
	taskstest.TestTaskRepository(t, tr)
}
