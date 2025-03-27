package sort

import (
	"os"
	"sort"
	"strings"
)

// sortKey returns a string for sorting:
// - For pseudo directories (".", ".."), it returns a key that ensures they come first.
// - For other hidden files (those starting with '.'), it returns the name without the leading dot.
// - For normal files, it returns the full lowercase name.
func SortKey(name string) string {
	if name == "." || name == ".." {
		// Prefix with spaces to ensure these come first.
		return "  " + name
	}
	if strings.HasPrefix(name, ".") {
		// Ignore the leading dot for sorting purposes.
		return strings.ToLower(name[1:])
	}
	return strings.ToLower(name)
}

// sortFiles orders file entries based on flags.
func SortFiles(files []os.DirEntry, flags map[string]bool) []os.DirEntry {
	// Default alphabetical sort based on our custom sort key.
	sort.Slice(files, func(i, j int) bool {
		return SortKey(files[i].Name()) < SortKey(files[j].Name())
	})

	// Sort by modification time if the "-t" flag is set.
	if flags["t"] {
		sort.Slice(files, func(i, j int) bool {
			infoI, _ := files[i].Info()
			infoJ, _ := files[j].Info()
			return infoI.ModTime().After(infoJ.ModTime())
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
