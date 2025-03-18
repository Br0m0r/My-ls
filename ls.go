package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Run is the entry point called from main.go.
func Run(args []string) {
	// Parse command-line options.
	opts := ParseArgs(args)
	// Set up output based on capture flag.
	out, cleanup := NewOutput(opts.Capture)
	// Run the core listing logic.
	runInternal(opts, out)
	// Call cleanup to close the capture file if necessary.
	cleanup()
}

// runInternal performs the core listing logic.
func runInternal(opts Options, out io.Writer) {
	paths := opts.Paths
	multiple := len(paths) > 1
	flagMap := opts.ToMap() // convert Options to map[string]bool for legacy functions

	for i, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(out, "Error: No such file or directory: %s\n", path)
			} else if os.IsPermission(err) {
				fmt.Fprintf(out, "Error: Permission denied: %s\n", path)
			} else {
				fmt.Fprintf(out, "Error: %v\n", err)
			}
			continue
		}

		if multiple {
			fmt.Fprintf(out, "%s:\n", path)
		}

		switch {
		case !info.IsDir():
			if opts.Long {
				pseudo := NewPseudoDirEntry(info, path)
				// Pass capture flag to displayLongFormat
				displayLongFormat([]fs.DirEntry{pseudo}, ".", out, opts.Capture)
			} else {
				fmt.Fprintf(out, "%s\n", path)
			}
		case opts.Recursive:
			// Pass capture flag to recursiveList
			recursiveList(path, flagMap, opts.Capture, out)
		default:
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Fprintf(out, "Error: %v\n", err)
				continue
			}
			files = filterFiles(files, flagMap, path)
			files = sortFiles(files, flagMap)
			// Pass capture flag to displayFiles
			displayFiles(files, path, flagMap, out, opts.Capture)
		}

		if i < len(paths)-1 {
			fmt.Fprintln(out)
		}
	}
}
