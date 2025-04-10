package filter

import (
	"io/fs"         
	"os"            
	"path/filepath" 
)

type pseudoDirectoryEntry struct {
	entryName string      // The name of the pseudo entry (e.g., "." or "..").
	fileInfo  os.FileInfo // The file information for this pseudo entry.
}

//Returns the name of the pseudo directory entry
func (p *pseudoDirectoryEntry) Name() string {
	return p.entryName
}

//Returns true if the pseudo entry represents a directory.
func (p *pseudoDirectoryEntry) IsDir() bool {
	return p.fileInfo.IsDir()
}

//Returns the file mode type from the pseudo entry's file info.
func (p *pseudoDirectoryEntry) Type() fs.FileMode {
	return p.fileInfo.Mode().Type()
}

//Returns the file information for the pseudo entry along with a nil error.
func (p *pseudoDirectoryEntry) Info() (fs.FileInfo, error) {
	return p.fileInfo, nil
}

//Creates a pseudoDirEntry for a given path and name using os.Stat.
func NewPseudoDirEntry(targetPath string, entryName string) (fs.DirEntry, error) {
	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, err 
	}
	
	return &pseudoDirectoryEntry{entryName: entryName, fileInfo: info}, nil
}

// Returns the parent directory of the given directory.
func GetParentDir(directoryPath string) string {
	if directoryPath == "." {
		// If current directory, get the working directory and then its parent.
		if wd, err := os.Getwd(); err == nil {
			return filepath.Dir(wd)
		}
	}
	// For any other directory, simply return its parent 
	return filepath.Dir(directoryPath)
}

// Filters directory entries based on the "-a" flag.
func FilterFiles(dirEntries []fs.DirEntry, flags map[string]bool, directoryPath string) []fs.DirEntry {
	if flags["a"] {
		var pseudoEntries []fs.DirEntry
		// Create a pseudo entry for the current directory "."
		if dotEntry, err := NewPseudoDirEntry(directoryPath, "."); err == nil {
			pseudoEntries = append(pseudoEntries, dotEntry)
		}
		// Get the parent directory and create a pseudo entry for ".."
		parentDirectory := GetParentDir(directoryPath)
		if dotDotEntry, err := NewPseudoDirEntry(parentDirectory, ".."); err == nil {
			pseudoEntries = append(pseudoEntries, dotDotEntry)
		}
		// Prepend the pseudo entries to the actual file list.
		return append(pseudoEntries, dirEntries...)
	}

	// If the "-a" flag is not set, filter out hidden files.
	var visibleEntries []fs.DirEntry
	for _, entry := range dirEntries {
		// Exclude files whose names start with a dot.
		if len(entry.Name()) > 0 && entry.Name()[0] != '.' {
			visibleEntries = append(visibleEntries, entry)
		}
	}
	return visibleEntries
}
