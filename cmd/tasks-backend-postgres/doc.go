// Package main combines the ginserver package with the gormrepo
// package and the GORM driver for postgres. It serves the Tasks
// API, and persists data in a Postgres database.
//
// It requires the following command-line flags OR environment
// variables.
//
// FLAG, ENVVAR, DEFAULT VALUE, DESCRIPTION
// -s, DB_SERVER,"", Name or IP address of a Postgres server
// -t, DB_SERVERPORT,"5432", Port on the Postgres server
// -u, DB_USERNAME, "", User ID
// -p, DB_PASSWORD, "", Password
// -d, DB_DATABASE, "", Name of database to connect to
//
// Each value above can also be read from a file, whose
// location can be specified using an environment variable that
// is named using a "FILE" suffix. E.g.:
//
//	DB_PASSWORDFILE="/etc/backendpassword.txt"
package main
