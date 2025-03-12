package main

import (
	"io/fs"
	"os"
	"path/filepath"
)

// pseudoDirEntry implements fs.DirEntry for pseudo entries like "." and "..".
type pseudoDirEntry struct {
	name string
	info os.FileInfo
}

func (p *pseudoDirEntry) Name() string {
	return p.name
}

func (p *pseudoDirEntry) IsDir() bool {
	return p.info.IsDir()
}

func (p *pseudoDirEntry) Type() fs.FileMode {
	return p.info.Mode().Type()
}

func (p *pseudoDirEntry) Info() (fs.FileInfo, error) {
	return p.info, nil
}

// newPseudoDirEntry creates a pseudoDirEntry for a given path and name.
func newPseudoDirEntry(path string, name string) (fs.DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &pseudoDirEntry{name: name, info: info}, nil
}

// getParentDir returns the parent directory of the given directory.
func getParentDir(dir string) string {
	return filepath.Dir(dir)
}

// filterFiles applies filters based on flags.
// If -a is set, it prepends pseudo entries for "." and ".." to the listing.
func filterFiles(files []fs.DirEntry, flags map[string]bool, dir string) []fs.DirEntry {
	if flags["a"] {
		var pseudoEntries []fs.DirEntry

		// Create pseudo entry for "."
		if dotEntry, err := newPseudoDirEntry(dir, "."); err == nil {
			pseudoEntries = append(pseudoEntries, dotEntry)
		}
		// Create pseudo entry for ".."
		parentDir := getParentDir(dir)
		if dotDotEntry, err := newPseudoDirEntry(parentDir, ".."); err == nil {
			pseudoEntries = append(pseudoEntries, dotDotEntry)
		}
		// Append the actual files (unfiltered) since -a should show everything.
		return append(pseudoEntries, files...)
	}

	// Without -a, filter out files that start with a dot.
	var visibleFiles []fs.DirEntry
	for _, file := range files {
		if len(file.Name()) > 0 && file.Name()[0] != '.' {
			visibleFiles = append(visibleFiles, file)
		}
	}
	return visibleFiles
}
