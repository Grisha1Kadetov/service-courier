package testutils

import (
	"os"
	"path/filepath"
)

func FindModuleRoot(dir string) (roots string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}
