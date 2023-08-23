package ginserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks/taskstest"
	"gorm.io/gorm"
)

type routeTestResult struct {
	statusCode int
	body       string
}

func testRoute(srv *Server, verb string, path string, body io.Reader) routeTestResult {
	req, _ := http.NewRequest(verb, path, body)
	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)

	responseData, _ := io.ReadAll(w.Body)
	return routeTestResult{
		statusCode: w.Code,
		body:       string(responseData),
	}
}

func TestIt(t *testing.T) {
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

	r := New(tr)

	result := testRoute(r, "GET", "/tasks", nil)
	if result.statusCode != http.StatusOK {
		t.Errorf("GET / failed. Response was:%v", result.body)
	}
}
