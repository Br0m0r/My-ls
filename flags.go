package main

import (
	"fmt"
	"os"
)

// Options holds parsed flag values and file/directory paths provided by the user.
type Options struct {
	Long      bool     // Use long listing format (-l)
	Recursive bool     // List subdirectories recursively (-R)
	ShowAll   bool     // Include entries starting with a dot (-a)
	TimeSort  bool     // Sort files by modification time (-t)
	Reverse   bool     // Reverse order while sorting (-r)
	Capture   bool     // Capture output to a file (-c)
	Paths     []string // File or directory paths provided as arguments
}

// ParseArgs processes command-line arguments and returns an Options struct populated with flags and paths.
func ParseArgs(args []string) Options {
	var opts Options

	// Loop through each argument provided in the command line.
	for _, arg := range args {
		// Check if the argument is a flag (starts with '-').
		if len(arg) > 0 && arg[0] == '-' {
			// Process each character in the flag argument (e.g., "-la" sets both l and a).
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					opts.Long = true // Enable long listing format.
				case 'R':
					opts.Recursive = true // Enable recursive directory listing.
				case 'a':
					opts.ShowAll = true // Include hidden files (starting with '.').
				case 't':
					opts.TimeSort = true // Sort files by modification time.
				case 'r':
					opts.Reverse = true // Reverse the sorting order.
				case 'c':
					opts.Capture = true // Capture the output to a file.
				case 'h':
					printUsage() // Print usage instructions.
					os.Exit(0)   // Exit the program after displaying help.
				default:
					// If an unknown flag is encountered, display an error message and usage.
					fmt.Printf("Unknown flag: -%c\n", ch)
					printUsage()
					os.Exit(1) // Exit with an error code.
				}
			}
		} else {
			// Non-flag arguments are considered as file or directory paths.
			opts.Paths = append(opts.Paths, arg)
		}
	}

	// If no paths were provided, default to the current directory.
	if len(opts.Paths) == 0 {
		opts.Paths = append(opts.Paths, ".")
	}
	return opts
}

// ToMap converts the Options struct to a map[string]bool for backward compatibility with functions
// that expect flags as a map. This excludes Capture and Paths as they are handled separately.
func (o Options) ToMap() map[string]bool {
	return map[string]bool{
		"l": o.Long,      // Long listing format flag.
		"R": o.Recursive, // Recursive listing flag.
		"a": o.ShowAll,   // Show all files (including hidden).
		"t": o.TimeSort,  // Sort by modification time.
		"r": o.Reverse,   // Reverse sort order.
	}
}

// printUsage displays the usage information for the program.
func printUsage() {
	fmt.Println("Usage: run.go [options] [path...]")
	fmt.Println("Options:")
	fmt.Println("  -l   Use long listing format")
	fmt.Println("  -R   List subdirectories recursively")
	fmt.Println("  -a   Include directory entries whose names begin with a dot (.)")
	fmt.Println("  -t   Sort by modification time, newest first")
	fmt.Println("  -r   Reverse order while sorting")
	fmt.Println("  -c   Capture output to file")
	fmt.Println("  -h   Display this help and exit")
}
