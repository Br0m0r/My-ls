package main

import (
	"fmt"
<<<<<<< HEAD
	"io/fs"
=======
	"io"
>>>>>>> 98c6a9a (fixed some formatting issues)
	"os"
)

// Run is the entry point called from main.go.
func Run(args []string) {
	capture := false
	// Remove "--capture" from arguments if present.
	newArgs := []string{}
	for _, arg := range args {
		if arg == "--capture" {
			capture = true
		} else {
			newArgs = append(newArgs, arg)
		}
	}

	// Set up the output writer as an io.Writer.
	var out io.Writer = os.Stdout
	if capture {
		f, err := os.Create("myls_output.txt")
		if err != nil {
			fmt.Println("Error creating capture file:", err)
			os.Exit(1)
		}
		defer f.Close()
		// out becomes a multiwriter so output goes both to stdout and the file.
		out = io.MultiWriter(os.Stdout, f)
	}

	runInternal(newArgs, out)
}

// runInternal does the actual processing.
func runInternal(args []string, out io.Writer) {
	flags, paths := parseFlags(args)
	// Default to current directory if no paths provided.
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	multiple := len(paths) > 1

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

<<<<<<< HEAD
		if !info.IsDir() {
			// Instead of just printing the file name,
			// create a pseudo DirEntry for the file.
			pseudoEntry, err := newPseudoDirEntry(path, path)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				// If -l flag is set, display in long format.
				if flags["l"] {
					// Since displayLongFormat expects a slice, wrap pseudoEntry in a slice.
					displayLongFormat([]fs.DirEntry{pseudoEntry}, ".")
				} else {
					// Otherwise, display just the file name.
					fmt.Println(path)
				}
			}
=======
		// If it's a file, simply print its name.
		if !info.IsDir() {
			fmt.Fprintf(out, "%s\n", path)
>>>>>>> 98c6a9a (fixed some formatting issues)
		} else {
			// Handle directories as before.
			if flags["R"] {
				recursiveList(path, flags, out)
			} else {
				files, err := os.ReadDir(path)
				if err != nil {
					fmt.Fprintf(out, "Error: %v\n", err)
					continue
				}
				files = filterFiles(files, flags, path)
				files = sortFiles(files, flags)
				displayFiles(files, path, flags, out)
			}
		}

		if i < len(paths)-1 {
			fmt.Fprintln(out)
		}
	}
}

// parseFlags is a simple implementation to separate flags from paths.
func parseFlags(args []string) (map[string]bool, []string) {
	flags := make(map[string]bool)
	var paths []string
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			for _, ch := range arg[1:] {
				flags[string(ch)] = true
			}
		} else {
			paths = append(paths, arg)
		}
	}
	return flags, paths
}
