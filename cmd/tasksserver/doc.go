// Package main combines the ginserver package with the gormrepo
// package. It serves the Tasks API, and static content from a
// subdirectory called wwwroot, under the current dierctory.
//
// By default, it listens on port 8080. This can be changed via
// a command line option called `-port`, or an environment
// variable called TASKS_PORT.
//
// Data can be stored in either SQLite or Postgres. This must be
// specified through a command line option called `-storage`, or
// an environment variable called TASKS_STORAGE. Possible values
// are "sqlite" and "postgres". The default is "sqlite".
//
// The "sqlite" option caused data to be stored in SQLite.
// It creates or uses a directory named `data` immediately under
// the current directory to store the database file, which is
// always called tasks.db.
//
// The "postgres" option causes data to be stored in Postgres.
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
