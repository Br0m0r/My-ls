package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// recursiveList handles `-R` flag and lists directories recursively.
func recursiveList(dir string, flags map[string]bool) {
	fmt.Println("\n" + dir + ":")

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	files = filterFiles(files, flags, dir)
	files = sortFiles(files, flags)
	displayFiles(files, dir, flags)

	// Iterate over directory entries.
	for _, file := range files {
		// Skip pseudo entries to avoid infinite recursion.
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		// Skip recursing into symbolic link directories.
		if file.Type()&os.ModeSymlink != 0 {
			continue
		}
		if info.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			recursiveList(subDir, flags)
		}
	}
}
