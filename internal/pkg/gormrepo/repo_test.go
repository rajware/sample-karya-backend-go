package gormrepo_test

import (
	"testing"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks/taskstest"
)

func TestGORMTaskRepository(t *testing.T) {
	tr := gormrepo.New()
	taskstest.TestTaskRepository(t, tr)
}
