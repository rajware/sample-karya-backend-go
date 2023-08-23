// Package opts contains utilities to provide configuration options
// to applications.
package opts

import (
	"flag"
	"os"
)

func parseflags() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func readvaluefromfile(filename string) string {
	resultbytes, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(resultbytes)
}

// GetOption gets an the value for an option. The value can be, in descending order
// of preference:
// - provided as a command-line option (optflag)
//
// - provided as an environment variable (envvarname)
//
// - provided in a file whose name is contained in an environment variable called
// <envvarname>FILE
//
// - the default value
func GetOption(optflag *string, envvarname string, defaultvalue string) string {
	parseflags()

	result := *optflag

	if result == "" || result == defaultvalue {
		result = os.Getenv(envvarname)
	}

	if result == "" {
		filename := os.Getenv(envvarname + "FILE")
		if filename != "" {
			result = readvaluefromfile(filename)
		}
	}

	if result == "" {
		result = defaultvalue
	}

	return result
}
