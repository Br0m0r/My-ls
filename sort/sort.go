package sort

import (
	"io/fs"
	"sort"
	"strings"
)

// SortKey returns a string key for sorting file names.
func SortKey(name string) string {
	switch name {
	case "-":
		return "0"
	case ".":
		return "1"
	case "..":
		return "2"
	}
	if strings.HasPrefix(name, ".") {
		return "3" + strings.ToLower(name[1:])
	}
	return "4" + strings.ToLower(name)
}

// SortFiles orders file entries based on flags.
func SortFiles(dirEntries []fs.DirEntry, optionFlags map[string]bool) []fs.DirEntry {
	// If the "-t" flag is set, sort by modification time (newest first)
	// with a secondary alphabetical order for files with identical modification times.
	if optionFlags["t"] {
		sort.SliceStable(dirEntries, func(i, j int) bool {
			entryInfoI, _ := dirEntries[i].Info()
			entryInfoJ, _ := dirEntries[j].Info()
			// If mod times are equal, sort alphabetically.
			if entryInfoI.ModTime() == entryInfoJ.ModTime() {
				return SortKey(dirEntries[i].Name()) < SortKey(dirEntries[j].Name())
			}
			// Otherwise, sort by modification time (newest first).
			return entryInfoI.ModTime().After(entryInfoJ.ModTime())
		})
	} else {
		// Default alphabetical sort based on our custom sort key.
		sort.Slice(dirEntries, func(i, j int) bool {
			return SortKey(dirEntries[i].Name()) < SortKey(dirEntries[j].Name())
		})
	}

	// Reverse order if the "-r" flag is set.
	if optionFlags["r"] {
		for i, j := 0, len(dirEntries)-1; i < j; i, j = i+1, j-1 {
			dirEntries[i], dirEntries[j] = dirEntries[j], dirEntries[i]
		}
	}

	return dirEntries
}
