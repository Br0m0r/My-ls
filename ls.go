package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Run is the entry point called from main.go.
// It parses arguments, sets up output, delegates to runInternal, and cleans up resources.
func Run(args []string) {
	opts := ParseArgs(args)
	out, cleanup := NewOutput(opts.Capture)
	runInternal(opts, out)
	cleanup()
}

// runInternal processes each path (file or directory) using the provided flags.
func runInternal(opts Options, out io.Writer) {
	paths := opts.Paths
	multiple := len(paths) > 1
	flagMap := opts.ToMap()

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
				displayLongFormat([]fs.DirEntry{pseudo}, ".", out, opts.Capture)
			} else {
				fmt.Fprintf(out, "%s\n", path)
			}
		case opts.Recursive:
			recursiveList(path, flagMap, opts.Capture, out)
		default:
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Fprintf(out, "Error: %v\n", err)
				continue
			}
			files = filterFiles(files, flagMap, path)
			files = sortFiles(files, flagMap)
			displayFiles(files, path, flagMap, out, opts.Capture)
		}

		if i < len(paths)-1 {
			fmt.Fprintln(out)
		}
	}
}
