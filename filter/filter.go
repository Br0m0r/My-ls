package filter

import (
	"io/fs"         // Provides file system interfaces like fs.DirEntry.
	"os"            // Used for file operations, e.g., os.Stat, os.Getwd.
	"path/filepath" // Provides functions for manipulating file paths.
)

// ----------------------------------------------------------
// Type: pseudoDirEntry
// Purpose: Implements fs.DirEntry for pseudo entries like "." and "..".
// ----------------------------------------------------------
type pseudoDirEntry struct {
	name string      // The name of the pseudo entry (e.g., "." or "..").
	info os.FileInfo // The file information for this pseudo entry.
}

// ----------------------------------------------------------
// Method: Name()
// Purpose: Returns the name of the pseudo directory entry.
// ----------------------------------------------------------
func (p *pseudoDirEntry) Name() string {
	return p.name
}

// ----------------------------------------------------------
// Method: IsDir()
// Purpose: Returns true if the pseudo entry represents a directory.
// ----------------------------------------------------------
func (p *pseudoDirEntry) IsDir() bool {
	return p.info.IsDir()
}

// ----------------------------------------------------------
// Method: Type()
// Purpose: Returns the file mode type from the pseudo entry's file info.
// ----------------------------------------------------------
func (p *pseudoDirEntry) Type() fs.FileMode {
	return p.info.Mode().Type()
}

// ----------------------------------------------------------
// Method: Info()
// Purpose: Returns the file information for the pseudo entry along with a nil error.
// ----------------------------------------------------------
func (p *pseudoDirEntry) Info() (fs.FileInfo, error) {
	return p.info, nil
}

// ----------------------------------------------------------
// Function: NewPseudoDirEntry
// Purpose: Creates a pseudoDirEntry for a given path and name using os.Stat.
// Parameters:
//   - path: The actual file system path to get FileInfo.
//   - name: The pseudo entry name (usually "." or "..").
//
// Returns:
//   - fs.DirEntry representing the pseudo entry, or an error if os.Stat fails.
//
// ----------------------------------------------------------
func NewPseudoDirEntry(path string, name string) (fs.DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err // Return error if the file/directory info cannot be obtained.
	}
	// Create and return a new pseudoDirEntry with the provided name and file info.
	return &pseudoDirEntry{name: name, info: info}, nil
}

// ----------------------------------------------------------
// Function: GetParentDir
// Purpose: Returns the parent directory of the given directory.
// Special Case: If the directory is ".", it resolves the absolute path first.
// ----------------------------------------------------------
func GetParentDir(dir string) string {
	if dir == "." {
		// If current directory, get the working directory and then its parent.
		if wd, err := os.Getwd(); err == nil {
			return filepath.Dir(wd)
		}
	}
	// For any other directory, simply return its parent using filepath.Dir.
	return filepath.Dir(dir)
}

// ----------------------------------------------------------
// Function: FilterFiles
// Purpose: Filters directory entries based on the "-a" flag.
// Behavior:
//   - If the "-a" flag is set, it prepends pseudo entries for "." and ".." to the file list.
//   - Otherwise, it filters out hidden files (those starting with a dot).
//
// ----------------------------------------------------------
func FilterFiles(files []fs.DirEntry, flags map[string]bool, dir string) []fs.DirEntry {
	if flags["a"] {
		var pseudoEntries []fs.DirEntry
		// Create a pseudo entry for the current directory "."
		if dotEntry, err := NewPseudoDirEntry(dir, "."); err == nil {
			pseudoEntries = append(pseudoEntries, dotEntry)
		}
		// Get the parent directory and create a pseudo entry for ".."
		parentDir := GetParentDir(dir)
		if dotDotEntry, err := NewPseudoDirEntry(parentDir, ".."); err == nil {
			pseudoEntries = append(pseudoEntries, dotDotEntry)
		}
		// Prepend the pseudo entries to the actual file list.
		return append(pseudoEntries, files...)
	}

	// If the "-a" flag is not set, filter out hidden files.
	var visibleFiles []fs.DirEntry
	for _, file := range files {
		// Exclude files whose names start with a dot.
		if len(file.Name()) > 0 && file.Name()[0] != '.' {
			visibleFiles = append(visibleFiles, file)
		}
	}
	return visibleFiles
}
