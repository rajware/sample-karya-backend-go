package ginserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
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
	// Set up data subdirectory
	datadir, err := taskstest.SetupDataDirectory()
	if err != nil {
		t.Logf("data directory creation failed with:%v", err)
		t.FailNow()
	}
	defer taskstest.RemoveDataDirectory()

	// Set up SQLite repo
	datafile := filepath.Join(datadir, "test.db")
	t.Logf("Opening data file %v", datafile)
	d := sqlite.Open(datafile)
	tr := gormrepo.New(d, &gorm.Config{})

	// Set up test server
	testserver := New(tr)

	// Set up file for static serving
	os.WriteFile("myfile.test", []byte("Hello3"), 0755)
	defer os.Remove("myfile.test")

	// Test static file
	result := testRoute(testserver, "GET", "/myfile.test", nil)
	if result.statusCode != http.StatusOK {
		t.Errorf("GET / failed. Response was:%v", result.body)
	}
	t.Logf("Result:%+v", result)

	// Test list of tasks
	result = testRoute(testserver, "GET", "/tasks", nil)
	if result.statusCode != http.StatusOK {
		t.Errorf("GET /tasks failed. Response was:%v", result.body)
	}

	if result.body != `{"data":[],"error":0,"message":"success"}` {
		t.Errorf("GET /tasks failed. Response was:%v", result.body)
	}
}
