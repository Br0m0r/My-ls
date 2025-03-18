package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

// pseudoDirEntry implements the fs.DirEntry interface for pseudo entries
// such as "." (current directory) and ".." (parent directory).
type pseudoDirEntry struct {
	name string      // The name of the pseudo entry, e.g., "." or ".."
	info os.FileInfo // File information for the pseudo entry
}

// Name returns the name of the pseudo directory entry.
func (p *pseudoDirEntry) Name() string {
	return p.name
}

// IsDir reports whether the pseudo entry represents a directory.
func (p *pseudoDirEntry) IsDir() bool {
	return p.info.IsDir()
}

// Type returns the mode type bits of the pseudo entry.
func (p *pseudoDirEntry) Type() fs.FileMode {
	return p.info.Mode().Type()
}

// Info returns the FileInfo associated with the pseudo entry.
func (p *pseudoDirEntry) Info() (fs.FileInfo, error) {
	return p.info, nil
}

// newPseudoDirEntry creates a new pseudoDirEntry for a given path and name.
// It uses os.Stat to obtain the FileInfo for the given path.
func newPseudoDirEntry(path string, name string) (fs.DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err // Return error if the file information can't be retrieved.
	}
	return &pseudoDirEntry{name: name, info: info}, nil
}

// getParentDir returns the parent directory of the given directory using filepath.Dir.
func getParentDir(dir string) string {
	return filepath.Dir(dir)
}

// filterFiles filters the provided directory entries based on the flags.
// When the "-a" (all) flag is set, it includes pseudo entries for "." and ".."
// along with the actual file entries. If "-a" is not set, it filters out hidden files.
func filterFiles(files []fs.DirEntry, flags map[string]bool, dir string) []fs.DirEntry {
	// If the "-a" flag is enabled, include all files, including hidden ones.
	if flags["a"] {
		var pseudoEntries []fs.DirEntry

		// Create a pseudo entry for the current directory ".".
		if dotEntry, err := newPseudoDirEntry(dir, "."); err == nil {
			pseudoEntries = append(pseudoEntries, dotEntry)
		}
		// Create a pseudo entry for the parent directory "..".
		parentDir := getParentDir(dir)
		if dotDotEntry, err := newPseudoDirEntry(parentDir, ".."); err == nil {
			pseudoEntries = append(pseudoEntries, dotDotEntry)
		}
		// Append the actual file entries after the pseudo entries.
		return append(pseudoEntries, files...)
	}

	// If "-a" is not enabled, filter out hidden files (files that start with '.').
	var visibleFiles []fs.DirEntry
	for _, file := range files {
		if len(file.Name()) > 0 && file.Name()[0] != '.' {
			visibleFiles = append(visibleFiles, file)
		}
	}
	return visibleFiles
}
