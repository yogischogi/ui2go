// Package toolbox implements common utility functions.
//
// I did not want to call this "util", because "util" would
// almost certainly lead to naming conflicts.
package toolbox

import (
	"errors"
	"os"
	"path/filepath"
)

// FindResourcesDir tries to locate a subdirectory that resides
// within the GOPATH environment variable.
//
// If the functions finds the specified subdirectory in one of the
// GOPATH directories it returns the full path of the directory.
func FindResourcesDir(resDir string) (dir string, err error) {
	goPaths := filepath.SplitList(os.Getenv("GOPATH"))
	for _, path := range goPaths {
		resourcesDir := filepath.Join(path, resDir)
		_, err := os.Stat(resourcesDir)
		if err == nil {
			return resourcesDir, err
		}
	}
	return "", errors.New("Could not locate resources directory.")
}
