package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/ginserver"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/opts"
	"golang.org/x/exp/slog"
)

var (
	storageOpt      string
	portOpt         string
	dbservernameOpt string
	dbserverportOpt string
	usernameOpt     string
	passwordOpt     string
	databasenameOpt string
)

func main() {
	flag.StringVar(&storageOpt, "storage", "", "storage engine to use")
	flag.StringVar(&portOpt, "port", "8080", "tasks api server port")
	flag.StringVar(&dbservernameOpt, "s", "localhost", "database server name or IP")
	flag.StringVar(&dbserverportOpt, "t", "5432", "database server port")
	flag.StringVar(&usernameOpt, "u", "", "user id")
	flag.StringVar(&passwordOpt, "p", "", "password")
	flag.StringVar(&databasenameOpt, "d", "", "database name")

	var ok bool

	portOpt, ok = verifyoption(
		&portOpt,
		"TASKS_PORT",
		"8080",
		"Fatal error: server port not specified",
	)
	if !ok {
		os.Exit(1)
	}
	port, err := strconv.ParseInt(portOpt, 10, 32)
	if err != nil {
		slog.Error("Fatal error: invalid server port", "port", portOpt)
		os.Exit((1))
	}

	storageOpt, ok = verifyoption(
		&storageOpt,
		"TASKS_STORAGE",
		"sqlite",
		"Fatal error: storage option not specified",
	)
	if !ok {
		os.Exit(1)
	}

	tr, err := getStorage(storageOpt)
	if err != nil {
		slog.Error("Fatal error: could not set up storage.", "Error", err)
		os.Exit(1)
	}

	// Start Server
	srv, err := ginserver.New(tr, int(port))
	if err != nil {
		slog.Error("Fatal error: could not start server.", "Error", err)
	}
	slog.Info("Starting server...")
	srv.Run()
}

func verifyoption(opt *string, envvarname string, defaultvalue string, message string) (string, bool) {
	result := opts.GetOption(opt, envvarname, defaultvalue)
	if result == "" {
		slog.Error(message)
		return result, false
	}
	return result, true
}
