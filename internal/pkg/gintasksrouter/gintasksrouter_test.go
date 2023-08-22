package gintasksrouter_test

import (
	"testing"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gintasksrouter"
)

func TestIt(t *testing.T) {
	r := gintasksrouter.New()
	r.Run()
}
