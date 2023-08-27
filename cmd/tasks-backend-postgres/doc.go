// Package main combines the ginserver package with the gormrepo
// package and the GORM driver for postgres. It serves the Tasks
// API, and persists data in a Postgres database.
//
// It requires the following command-line flags OR environment
// variables.
//
// FLAG, ENVVAR, DEFAULT VALUE, DESCRIPTION
// -s, TASKS_DBSERVER,"", Name or IP address of a Postgres server
// -t, TASKS_DBPORT,"5432", Port on the Postgres server
// -u, TASKS_USERNAME, "", User ID
// -p, TASKS_PASSWORD, "", Password
// -d, TASKS_DATABASE, "", Name of database to connect to
//
// Each value above can also be read from a file, whose
// location can be specified using an environment variable that
// is named using a "FILE" suffix. E.g.:
//
//	DB_PASSWORDFILE="/etc/backendpassword.txt"
package main
