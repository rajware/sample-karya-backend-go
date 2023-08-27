package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rajware/sample-tasks-backend-go/internal/pkg/ginserver"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/opts"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func verifyoption(opt *string, envvarname string, defaultvalue string, message string) (string, bool) {
	result := opts.GetOption(opt, envvarname, defaultvalue)
	if result == "" {
		slog.Error(message)
		return result, false
	}
	return result, true
}

func main() {
	var (
		servernameOpt   string
		serverportOpt   string
		usernameOpt     string
		passwordOpt     string
		databasenameOpt string
	)

	flag.StringVar(&servernameOpt, "s", "", "database server name or IP")
	flag.StringVar(&serverportOpt, "t", "5432", "database server port")
	flag.StringVar(&usernameOpt, "u", "", "user id")
	flag.StringVar(&passwordOpt, "p", "", "password")
	flag.StringVar(&databasenameOpt, "d", "", "database name")

	allOk := true

	servernameOpt, ok := verifyoption(&servernameOpt, "TASKS_DBSERVER", "", "Postgres server name not provided.")
	allOk = allOk && ok

	serverportOpt, ok = verifyoption(&serverportOpt, "TASKS_DBPORT", "5432", "Postgres server port not provided.")
	allOk = allOk && ok

	usernameOpt, ok = verifyoption(&usernameOpt, "TASKS_USERNAME", "", "user name not provided.")
	allOk = allOk && ok

	passwordOpt, ok = verifyoption(&passwordOpt, "TASKS_PASSWORD", "", "password not provided.")
	allOk = allOk && ok

	databasenameOpt, ok = verifyoption(&databasenameOpt, "TASKS_DBNAME", "", "database name not provided.")
	allOk = allOk && ok

	if !allOk {
		os.Exit(1)
	}

	pgdsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Calcutta",
		servernameOpt,
		usernameOpt,
		passwordOpt,
		databasenameOpt,
		serverportOpt,
	)
	d := postgres.Open(pgdsn)
	tr := gormrepo.New(d, &gorm.Config{})

	// Start Server
	srv, err := ginserver.New(tr, 8080)
	if err != nil {
		slog.Error("Fatal error: could not start server.", "Error", err)
	}
	slog.Info("Starting server...")
	srv.Run()
}
