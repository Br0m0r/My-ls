package sort

import (
	"io/fs"
	"sort"
	"strings"
)

// SortKey returns a string key for sorting file names.
// It assigns special keys to "-", ".", and ".." so that they sort in the order:
// "-" (first), then "." then "..", followed by all other names.
// For other names, if the name starts with a dot, the dot is ignored for sorting.
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
		// Remove the leading dot for sorting purposes.
		return "3" + strings.ToLower(name[1:])
	}
	return "3" + strings.ToLower(name)
}

// SortFiles orders file entries based on flags.
func SortFiles(files []fs.DirEntry, flags map[string]bool) []fs.DirEntry {
	// If the "-t" flag is set, sort by modification time (newest first)
	// with a secondary alphabetical order for files with identical modification times.
	if flags["t"] {
		sort.SliceStable(files, func(i, j int) bool {
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			// If mod times are equal, sort alphabetically.
			if infoI.ModTime() == infoJ.ModTime() {
				return SortKey(files[i].Name()) < SortKey(files[j].Name())
			}
			// Otherwise, sort by modification time (newest first).
			return infoI.ModTime().After(infoJ.ModTime())
		})
	} else {
		// Default alphabetical sort based on our custom sort key.
		sort.Slice(files, func(i, j int) bool {
			return SortKey(files[i].Name()) < SortKey(files[j].Name())
		})
	}

	// Reverse order if the "-r" flag is set.
	if flags["r"] {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	return files
}
