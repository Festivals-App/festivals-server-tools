package servertools

import (
	"os"
	"slices"
	"strings"
)

// Checks if the binary was executed with the --container="/path/to/container" argument and returns the specified path or an empty string.
func ContainerPathArgument() string {

	args := os.Args[1:]
	for i := range args {
		arg := args[i]
		values := strings.Split(arg, "=") // input is trusted
		if len(values) == 2 {
			cmd := values[0]
			value := values[1]
			if cmd == "--container" {
				return value
			}
		}
	}
	return ""
}

// Checks if the binary was executed with the --debug argument.
func DebugArgument() bool {
	return slices.Contains(os.Args[1:], "--debug")
}

/*
func IsRunningInDebug() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return isPresent
}

func IsRunningInProduction() bool {
	_, isPresent := os.LookupEnv("DEBUG")
	return !isPresent
}
*/
