package main

import (
	"os"
	"sort"
	"strings"
)

// sortFiles orders file entries based on flags.
func sortFiles(files []os.DirEntry, flags map[string]bool) []os.DirEntry {
	// Default alphabetical sort (case-insensitive).
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	// Sort by modification time if "-t" is set.
	if flags["t"] {
		sort.Slice(files, func(i, j int) bool {
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			return infoI.ModTime().After(infoJ.ModTime())
		})
	}

	// Reverse order if "-r" is set.
	if flags["r"] {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	return files
}
