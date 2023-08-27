package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/ginserver"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/tasks"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var storageMap = map[string]func() (tasks.TaskRepository, error){
	"sqlite":   getSqliteStorage,
	"postgres": getPostgresStorage,
}

func getStorage(storage string) (tasks.TaskRepository, error) {
	result, ok := storageMap[storage]
	if !ok {
		return nil, errors.New("invalid storage option '" + storage + "'")
	}
	return result()
}

func getSqliteStorage() (tasks.TaskRepository, error) {
	slog.Info("Setting up SQLite storage...")

	// Set up data directory
	datadir, err := ginserver.EnsureSubdirectory("data")
	if err != nil {
		return nil, errors.New("could not locate data directory:" + err.Error())
	}

	// Set up SQLite repository
	datafile := filepath.Join(datadir, "tasks.db")
	slog.Info("Opening data file", "file", datafile)
	d := sqlite.Open(datafile)
	tr := gormrepo.New(d, &gorm.Config{})

	slog.Info("SQLite storage set up.")
	return tr, nil
}

func getPostgresStorage() (tasks.TaskRepository, error) {
	slog.Info("Setting up Postgres storage...")
	allOk := true

	var ok bool
	dbservernameOpt, ok = verifyoption(&dbservernameOpt, "TASKS_DBSERVER", "db", "Postgres server name not provided.")
	allOk = allOk && ok

	dbserverportOpt, ok = verifyoption(&dbserverportOpt, "TASKS_DBPORT", "5432", "Postgres server port not provided.")
	allOk = allOk && ok

	usernameOpt, ok = verifyoption(&usernameOpt, "TASKS_USERNAME", "", "user name not provided.")
	allOk = allOk && ok

	passwordOpt, ok = verifyoption(&passwordOpt, "TASKS_PASSWORD", "", "password not provided.")
	allOk = allOk && ok

	databasenameOpt, ok = verifyoption(&databasenameOpt, "TASKS_DBNAME", "", "database name not provided.")
	allOk = allOk && ok

	if !allOk {
		return nil, errors.New("storage provider connection parameters not set")
	}

	pgdsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Calcutta",
		dbservernameOpt,
		usernameOpt,
		passwordOpt,
		databasenameOpt,
		dbserverportOpt,
	)
	d := postgres.Open(pgdsn)
	tr := gormrepo.New(d, &gorm.Config{})

	slog.Info("Postgres storage set up.")
	return tr, nil
}
