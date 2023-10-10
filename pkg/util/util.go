// TODO: get rid of this package
package util

import (
	"os"
	"path/filepath"
)

// Returns absolute executable path of the application
func ExecPath() string {
	path, err := os.Executable()
	if err != nil {
		// Try to get from args
		path = os.Args[0]
	}
	path, _ = filepath.EvalSymlinks(path)
	path, _ = filepath.Abs(path)
	return path
}

// Reports whether the file at path exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
