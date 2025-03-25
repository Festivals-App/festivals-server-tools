package servertools

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

// Checks if a file exists and isn't a directory.
// see: https://golangcode.com/check-if-a-file-exists/
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

// Checks if the binary was executed with the --container argument and returns the specified value or an empty string.
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

// ExpandTilde expands a leading tilde (~) to the user's home directory who runs the binary.
func ExpandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}

func IsEntitled(keys []string, endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api-Key"] != nil {

			key := r.Header["Api-Key"][0]

			if contains(keys, key) {
				endpoint(w, r)
			} else {
				RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

		} else {
			RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
	})
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
