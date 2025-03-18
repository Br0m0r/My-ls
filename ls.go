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
	out := NewOutput(opts.Capture)
	runInternal(opts, out)
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
			// For files, use long listing if -l flag is set.
			if opts.Long {
				pseudo := NewPseudoDirEntry(info, path)
				// Wrap the pseudo entry in a slice to reuse displayLongFormat.
				displayLongFormat([]fs.DirEntry{pseudo}, ".", out)
			} else {
				fmt.Fprintf(out, "%s\n", path)
			}
		case opts.Recursive:
			recursiveList(path, flagMap, out)
		default:
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Fprintf(out, "Error: %v\n", err)
				continue
			}
			files = filterFiles(files, flagMap, path)
			files = sortFiles(files, flagMap)
			displayFiles(files, path, flagMap, out)
		}

		if i < len(paths)-1 {
			fmt.Fprintln(out)
		}
	}
}
