package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// recursiveList lists directories recursively.
func recursiveList(dir string, flags map[string]bool, capture bool, out io.Writer) {
	fmt.Fprintf(out, "\n%s:\n", dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(out, "Error: %v\n", err)
		return
	}

	files = filterFiles(files, flags, dir)
	files = sortFiles(files, flags)
	displayFiles(files, dir, flags, out, capture)

	for _, file := range files {
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		if file.Type()&os.ModeSymlink != 0 {
			continue
		}
		if info.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			recursiveList(subDir, flags, capture, out)
		}
	}
}
