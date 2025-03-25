package servertools

import (
	"errors"
	"io/fs"
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

// ExpandTilde expands a leading tilde (~) to the user's home directory who runs the binary.
func ExpandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			return strings.Replace(path, "~", home, 1)
		}
	}
	return path
}
