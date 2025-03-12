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

	// Iterate over directory entries
	for _, file := range files {
		if file.IsDir() {
			subDir := filepath.Join(dir, file.Name())
			recursiveList(subDir, flags)
		}
	}
}
