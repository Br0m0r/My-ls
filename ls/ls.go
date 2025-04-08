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

// Run is the entry point called from cmd/myls/main.go.
// It parses arguments, sets up output, delegates to RunInternal, and cleans up resources.
func Run(args []string) {
	opts := flags.ParseArgs(args)
	out, cleanup := output.NewOutput(opts.Capture)
	RunInternal(opts, out)
	cleanup()
}

// RunInternal separates file and directory arguments.
// Files are processed first, and then directories.
// If the recursive flag (-R) is set, each directory is listed recursively.
func RunInternal(opts flags.Options, out io.Writer) {
	paths := opts.Paths
	flagMap := opts.ToMap()

	var filePaths []string
	var dirPaths []string

	// Separate file and directory arguments.
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			if strings.Contains(err.Error(), "not a directory") && strings.HasSuffix(path, "/") {
				fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': Not a directory\n", path)
			} else if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': No such file or directory\n", path)
			} else if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot open '%s': Permission denied\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
			}
			continue
		}
		if info.IsDir() {
			dirPaths = append(dirPaths, path)
		} else {
			filePaths = append(filePaths, path)
		}
	}

	// Process file arguments first.
	for _, path := range filePaths {
		info, _ := os.Lstat(path)
		if opts.Long {
			// For a single file, do not print the "total" line.
			pseudo := utils.NewPseudoDirEntry(info, path)
			display.DisplayLongFormat([]fs.DirEntry{pseudo}, ".", out, opts.Capture, false)
		} else {
			fmt.Fprintf(out, "%s\n", path)
		}
	}

	// If both file and directory arguments exist, print a newline as separator.
	if len(filePaths) > 0 && len(dirPaths) > 0 {
		fmt.Fprintln(out)
	}

	// Process directory arguments.
	// Print header (i.e. directory name + colon) if more than one directory or if files are listed above.
	multipleHeaders := (len(dirPaths) > 1) || (len(filePaths) > 0)
	for i, dir := range dirPaths {
		if multipleHeaders {
			fmt.Fprintf(out, "%s:\n", dir)
		}
		if opts.Recursive {
			// Recursive listing for directories.
			recursive.RecursiveList(dir, flagMap, opts.Capture, out)
		} else {
			entries, err := os.ReadDir(dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
				continue
			}
			entries = filter.FilterFiles(entries, flagMap, dir)
			entries = sort.SortFiles(entries, flagMap)
			display.DisplayFiles(entries, dir, flagMap, out, opts.Capture)
		}
		if i < len(dirPaths)-1 {
			fmt.Fprintln(out)
		}
	}
}
