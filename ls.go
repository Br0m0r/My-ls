package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// Run is the entry point called from main.go.
// It parses command-line arguments, sets up output (including capture mode),
// delegates to the core listing logic, and finally cleans up resources.
func Run(args []string) {
	// Parse command-line options into an Options struct.
	opts := ParseArgs(args)
	// Set up output writer based on whether output capture is enabled (-c flag).
	// NewOutput returns both an io.Writer (which could be stdout or a file) and a cleanup function.
	out, cleanup := NewOutput(opts.Capture)
	// Run the core file/directory listing logic.
	runInternal(opts, out)
	// Call the cleanup function to close the capture file if necessary.
	cleanup()
}

// runInternal performs the core file listing logic.
// It processes each path provided (could be a file or a directory), applying different
// display methods based on the flags set in the Options.
func runInternal(opts Options, out io.Writer) {
	// Get all paths provided by the user.
	paths := opts.Paths
	// Check if more than one path is provided so that headers can be printed.
	multiple := len(paths) > 1
	// Convert Options to a map[string]bool for compatibility with legacy functions.
	flagMap := opts.ToMap()

	// Iterate over each path provided.
	for i, path := range paths {
		// Retrieve file information for the current path.
		info, err := os.Stat(path)
		if err != nil {
			// Handle errors: path does not exist or permission issues.
			if os.IsNotExist(err) {
				fmt.Fprintf(out, "Error: No such file or directory: %s\n", path)
			} else if os.IsPermission(err) {
				fmt.Fprintf(out, "Error: Permission denied: %s\n", path)
			} else {
				fmt.Fprintf(out, "Error: %v\n", err)
			}
			// Skip to the next path if there's an error.
			continue
		}

		// If there are multiple paths, print a header with the current path.
		if multiple {
			fmt.Fprintf(out, "%s:\n", path)
		}

		// Determine how to display the current path based on its type and flags.
		switch {
		// Case: The path is a file (not a directory).
		case !info.IsDir():
			// If long listing format (-l) is requested, display detailed information.
			if opts.Long {
				// Create a pseudo directory entry for the file to be used by displayLongFormat.
				pseudo := NewPseudoDirEntry(info, path)
				// Display file details in long format (ls -l style).
				displayLongFormat([]fs.DirEntry{pseudo}, ".", out, opts.Capture)
			} else {
				// Otherwise, simply print the file name.
				fmt.Fprintf(out, "%s\n", path)
			}

		// Case: The path is a directory and recursive listing (-R) is enabled.
		case opts.Recursive:
			// List the directory recursively.
			recursiveList(path, flagMap, opts.Capture, out)

		// Case: Standard directory listing (non-recursive).
		default:
			// Read the directory entries.
			files, err := os.ReadDir(path)
			if err != nil {
				fmt.Fprintf(out, "Error: %v\n", err)
				continue
			}
			// Apply filtering based on flags (e.g., hide dot files unless -a is set).
			files = filterFiles(files, flagMap, path)
			// Sort the files based on the selected sorting options (alphabetical, time sort, reverse).
			files = sortFiles(files, flagMap)
			// Display the files. The displayFiles function handles whether to print in long format or simple mode.
			displayFiles(files, path, flagMap, out, opts.Capture)
		}

		// If there are multiple paths, separate outputs with a newline.
		if i < len(paths)-1 {
			fmt.Fprintln(out)
		}
	}
}
