package recursive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"eles/display"
	"eles/filter"
	"eles/sort"
)

// joinDisplayPath joins parent and child directory names, preserving the "./" prefix when the parent is "." or starts with "./".
func joinDisplayPath(parentPath, childName string) string {
	if parentPath == "." {
		return "./" + childName
	}
	if strings.HasPrefix(parentPath, "./") {
		if strings.HasSuffix(parentPath, "/") {
			return parentPath + childName
		}
		return parentPath + "/" + childName
	}
	return filepath.Join(parentPath, childName)
}

// RecursiveList lists directories recursively.
func RecursiveList(directoryPath string, optionFlags map[string]bool, captureOutput bool, outputWriter io.Writer) {
	fmt.Fprintf(outputWriter, "\n%s:\n", directoryPath)

	dirEntries, err := os.ReadDir(directoryPath)
	if err != nil {
		fmt.Fprintf(outputWriter, "Error: %v\n", err)
		return
	}

	dirEntries = filter.FilterFiles(dirEntries, optionFlags, directoryPath)
	dirEntries = sort.SortFiles(dirEntries, optionFlags)
	display.DisplayFiles(dirEntries, directoryPath, optionFlags, outputWriter, captureOutput)

	for _, entry := range dirEntries {
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			continue
		}
		if entryInfo.IsDir() {
			subDirectoryPath := joinDisplayPath(directoryPath, entry.Name())
			RecursiveList(subDirectoryPath, optionFlags, captureOutput, outputWriter)
		}
	}
}
