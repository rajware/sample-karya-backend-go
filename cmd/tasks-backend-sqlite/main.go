package main

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/ginserver"
	"github.com/rajware/sample-tasks-backend-go/internal/pkg/gormrepo"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func setupDataDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	datadir := filepath.Join(wd, "data")
	err = os.MkdirAll(datadir, 0755)
	if err != nil {
		return "", err
	}

	return datadir, nil
}

func main() {
	// Set up data directory
	datadir, err := setupDataDirectory()
	if err != nil {
		slog.Error("could not locate data directory")
		os.Exit(1)
	}

	// Set up SQLite repository
	datafile := filepath.Join(datadir, "tasks.db")
	slog.Info("Opening data file.", "file", datafile)
	d := sqlite.Open(datafile)
	tr := gormrepo.New(d, &gorm.Config{})

	// Start Server
	srv, err := ginserver.New(tr, 8080)
	if err != nil {
		slog.Error("Fatal error: could not start server.", "Error", err)
	}
	slog.Info("Starting server...")
	srv.Run()
}
