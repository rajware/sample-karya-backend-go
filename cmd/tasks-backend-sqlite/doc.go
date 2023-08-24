// Package main combines the ginserver package with the gormrepo
// package and the GORM driver for SQLite. It serves the Tasks
// API, and persists data in a SQLite database file.
//
// It creates or uses a directory named `data` immediately under
// the current directory to store the database file, which is
// always called tasks.db.
package main
