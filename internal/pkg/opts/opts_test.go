package opts_test

import (
	"flag"
	"os"
	"testing"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/opts"
)

var nameopt string

func init() {
	flag.StringVar(&nameopt, "name", "", "A name")
	os.Args = []string{"test", "-name", "Hello"}
}

func TestCommandLineOpt(t *testing.T) {
	os.Setenv("MY_NAME", "Hello2")
	os.WriteFile("myname.txt", []byte("Hello3"), 0755)
	defer os.Remove("myname.txt")
	os.Setenv("MY_NAMEFILE", "myname.txt")

	testOpt(t, "command line option", "Hello")

	nameopt = ""
	testOpt(t, "environment variable", "Hello2")

	os.Unsetenv("MY_NAME")
	testOpt(t, "file in environment variable", "Hello3")

	os.Unsetenv("MY_NAMEFILE")
	testOpt(t, "default value", "just a name")
}

func testOpt(t *testing.T, test string, expectedresult string) {
	t.Logf("Testing %v...", test)
	result := opts.GetOption(&nameopt, "MY_NAME", "just a name")
	if result != expectedresult {
		t.Errorf("%v not picked up. Result was:%v", test, result)
	}
}
