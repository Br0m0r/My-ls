package ls

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"eles/display"
	"eles/filter"
	"eles/flags"
	"eles/output"
	"eles/recursive"
	"eles/sort"
	"eles/utils"
)

// Run is the entry point called from cmd/myls/main.go.
// It parses arguments, sets up output, delegates to runInternal, and cleans up resources.
func Run(args []string) {
	opts := flags.ParseArgs(args)
	out, cleanup := output.NewOutput(opts.Capture)
	RunInternal(opts, out)
	cleanup()
}

func RunInternal(opts flags.Options, out io.Writer) {
	paths := opts.Paths
	multiple := len(paths) > 1
	flagMap := opts.ToMap()

	// First, loop over all arguments and print errors immediately to stderr.
	// Also collect the valid paths for later processing.
	var validPaths []string
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot access '%s': No such file or directory\n", path)
			} else if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "my-ls: cannot open directory '%s': Permission denied\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
			}
		} else {
			validPaths = append(validPaths, path)
		}
	}

	// Now, process and display the valid paths.
	for i, path := range validPaths {
		info, _ := os.Stat(path) // safe since we already checked the error
		if multiple {
			fmt.Fprintf(out, "%s:\n", path)
		}

		switch {
		// Non-directory: display file directly.
		case !info.IsDir():
			if opts.Long {
				pseudo := utils.NewPseudoDirEntry(info, path)
				display.DisplayLongFormat([]fs.DirEntry{pseudo}, ".", out, opts.Capture)
			} else {
				fmt.Fprintf(out, "%s\n", path)
			}
		// Recursive listing for directories if -R is set.
		case opts.Recursive:
			recursive.RecursiveList(path, flagMap, opts.Capture, out)
		// Regular directory listing.
		default:
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
				continue
			}
			files = filter.FilterFiles(files, flagMap, path)
			files = sort.SortFiles(files, flagMap)
			display.DisplayFiles(files, path, flagMap, out, opts.Capture)
		}

		if i < len(validPaths)-1 {
			fmt.Fprintln(out)
		}
	}
}
