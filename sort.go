package main

import (
	"os"
	"sort"
	"strings"
)

// sortFiles orders the slice of file entries based on the flags provided.
// It performs three types of sorting:
// 1. Alphabetical sort (default) in a case-insensitive manner.
// 2. If the "-t" flag is set, it sorts files by modification time (newest first).
// 3. If the "-r" flag is set, it reverses the sorted order.
func sortFiles(files []os.DirEntry, flags map[string]bool) []os.DirEntry {
	// Default alphabetical sorting (case-insensitive).
	// Using sort.Slice with an anonymous function that compares the file names.
	sort.Slice(files, func(i, j int) bool {
		// Convert both file names to lowercase to ensure case-insensitive comparison.
		return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
	})

	// If the "-t" flag is set, sort the files by their modification time.
	if flags["t"] {
		sort.Slice(files, func(i, j int) bool {
			// Retrieve file info for both files.
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			// Compare modification times: more recent files come first.
			return infoI.ModTime().After(infoJ.ModTime())
		})
	}

	// If the "-r" flag is set, reverse the order of the sorted slice.
	if flags["r"] {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			// Swap the elements at indices i and j.
			files[i], files[j] = files[j], files[i]
		}
	}

	// Return the sorted (and possibly reversed) slice of file entries.
	return files
}
