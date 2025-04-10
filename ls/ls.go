package ls

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"eles/display"
	"eles/filter"
	"eles/flags"
	"eles/output"
	"eles/recursive"
	"eles/sort"
	"eles/utils"
)

// Run is the entry point called from main.go.
func Run(arguments []string) {
	options := flags.ParseArgs(arguments)
	outputWriter, cleanupFunc := output.NewOutput(options.Capture)
	RunInternal(options, outputWriter)
	cleanupFunc()
}

// RunInternal separates file and directory arguments.
// Files are processed first, and then directories.
// If the recursive flag (-R) is set, each directory is listed recursively.
func RunInternal(options flags.Options, outputWriter io.Writer) {
	inputPaths := options.Paths
	optionFlags := options.ToMap()

	var fileArgumentPaths []string
	var directoryArgumentPaths []string

	// Separate file and directory arguments.
	for _, currentPath := range inputPaths {
		fileInfo, err := os.Lstat(currentPath)
		if err != nil {
			if strings.Contains(err.Error(), "not a directory") && strings.HasSuffix(currentPath, "/") {
				fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': Not a directory\n", currentPath)
			} else if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': No such file or directory\n", currentPath)
			} else if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot open '%s': Permission denied\n", currentPath)
			} else {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
			}
			continue
		}
		if fileInfo.IsDir() {
			directoryArgumentPaths = append(directoryArgumentPaths, currentPath)
		} else {
			fileArgumentPaths = append(fileArgumentPaths, currentPath)
		}
	}

	// Process file arguments first.
	for _, filePath := range fileArgumentPaths {
		fileInfo, _ := os.Lstat(filePath)
		if options.Long {
			// For a single file, do not print the "total" line.
			pseudoEntry := utils.NewPseudoDirEntry(fileInfo, filePath)
			display.DisplayLongFormat([]fs.DirEntry{pseudoEntry}, ".", outputWriter, options.Capture, false)
		} else {
			fmt.Fprintf(outputWriter, "%s\n", filePath)
		}
	}

	// If both file and directory arguments exist, print a newline as separator.
	if len(fileArgumentPaths) > 0 && len(directoryArgumentPaths) > 0 {
		fmt.Fprintln(outputWriter)
	}

	// Process directory arguments.
	// Print header (i.e. directory name + colon) if more than one directory or if files are listed above.
	multipleHeaders := (len(directoryArgumentPaths) > 1) || (len(fileArgumentPaths) > 0)
	for index, directoryPath := range directoryArgumentPaths {
		if multipleHeaders {
			fmt.Fprintf(outputWriter, "%s:\n", directoryPath)
		}
		if options.Recursive {
			// Recursive listing for directories.
			recursive.RecursiveList(directoryPath, optionFlags, options.Capture, outputWriter)
		} else {
			directoryEntries, err := os.ReadDir(directoryPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
				continue
			}
			directoryEntries = filter.FilterFiles(directoryEntries, optionFlags, directoryPath)
			directoryEntries = sort.SortFiles(directoryEntries, optionFlags)
			display.DisplayFiles(directoryEntries, directoryPath, optionFlags, outputWriter, options.Capture)
		}
		if index < len(directoryArgumentPaths)-1 {
			fmt.Fprintln(outputWriter)
		}
	}
}
