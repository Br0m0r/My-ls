package filter

import (
	"io/fs"
	"os"
	"path/filepath"
)

// pseudoDirEntry implements fs.DirEntry for pseudo entries like "." and "..".
type pseudoDirEntry struct {
	name string      // e.g., "." or ".."
	info os.FileInfo // File info for the pseudo entry
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

// newPseudoDirEntry creates a pseudoDirEntry using os.Stat.
func NewPseudoDirEntry(path string, name string) (fs.DirEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &pseudoDirEntry{name: name, info: info}, nil
}

// getParentDir returns the parent directory of the given directory.
// If the directory is ".", it resolves the absolute path to determine its actual parent.
func GetParentDir(dir string) string {
	if dir == "." {
		if wd, err := os.Getwd(); err == nil {
			return filepath.Dir(wd)
		}
	}
	return filepath.Dir(dir)
}

// filterFiles filters directory entries based on the "-a" flag.
// If "-a" is set, it prepends pseudo entries for "." and "..".
func FilterFiles(files []fs.DirEntry, flags map[string]bool, dir string) []fs.DirEntry {
	if flags["a"] {
		var pseudoEntries []fs.DirEntry
		if dotEntry, err := NewPseudoDirEntry(dir, "."); err == nil {
			pseudoEntries = append(pseudoEntries, dotEntry)
		}
		parentDir := GetParentDir(dir)
		if dotDotEntry, err := NewPseudoDirEntry(parentDir, ".."); err == nil {
			pseudoEntries = append(pseudoEntries, dotDotEntry)
		}
		return append(pseudoEntries, files...)
	}

	var visibleFiles []fs.DirEntry
	for _, file := range files {
		if len(file.Name()) > 0 && file.Name()[0] != '.' {
			visibleFiles = append(visibleFiles, file)
		}
	}
	return visibleFiles
}
