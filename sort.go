package main

import (
	"os"
	"sort"
	"strings"
)

// sortFiles applies sorting rules based on flags.
func sortFiles(files []os.DirEntry, flags map[string]bool) []os.DirEntry {
	// Sort alphabetically by default
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	// Sort by modification time if -t is used
	if flags["t"] {
		sort.Slice(files, func(i, j int) bool {
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			return infoI.ModTime().After(infoJ.ModTime())
		})
	}

	// Reverse order if -r is used
	if flags["r"] {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	return files
}
