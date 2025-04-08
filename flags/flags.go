package flags

import (
	"fmt"
	"os"
)

// Options holds parsed flag values and file/directory paths.
type Options struct {
	Long      bool     // Use long listing format (-l)
	Recursive bool     // List directories recursively (-R)
	ShowAll   bool     // Include hidden entries (-a)
	TimeSort  bool     // Sort by modification time (-t)
	Reverse   bool     // Reverse sort order (-r)
	Capture   bool     // Capture output to a file (-c)
	Paths     []string // File or directory paths provided as arguments
}

func ParseArgs(args []string) Options {
	var opts Options
	endOfOptions := false

	for _, arg := range args {
		// if we haven't seen "--", then any arg starting with '-' is options
		if !endOfOptions && len(arg) > 0 && arg[0] == '-' {
			// Special case: if the argument is exactly "--", then subsequent args are literal file names.
			if arg == "--" {
				endOfOptions = true
				continue
			}
			// Otherwise, parse each flag character in this argument.
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					opts.Long = true
				case 'R':
					opts.Recursive = true
				case 'a':
					opts.ShowAll = true
				case 't':
					opts.TimeSort = true
				case 'r':
					opts.Reverse = true
				case 'c':
					opts.Capture = true
				case 'h':
					printUsage()
					os.Exit(0)
				default:
					fmt.Printf("Unknown flag: -%c\n", ch)
					printUsage()
					os.Exit(1)
				}
			}
			// Move on to the next argument.
			continue
		}
		// Either endOfOptions is true, or the argument doesn't start with '-', so treat as a path.
		opts.Paths = append(opts.Paths, arg)
	}

	// Default behavior: if no paths provided, use current directory.
	if len(opts.Paths) == 0 {
		opts.Paths = append(opts.Paths, ".")
	}
	return opts
}

// ToMap converts Options to a map for compatibility with other functions.
func (o Options) ToMap() map[string]bool {
	// ----------------------------------------------------------
	// Create a Map with Flag Names as Keys and Their Boolean Values
	// ----------------------------------------------------------
	return map[string]bool{
		"l": o.Long,
		"R": o.Recursive,
		"a": o.ShowAll,
		"t": o.TimeSort,
		"r": o.Reverse,
	}
}

// printUsage prints a help message describing how to use the command.
func printUsage() {
	// ----------------------------------------------------------
	// Print General Usage Information
	// ----------------------------------------------------------
	fmt.Println("Usage: myls [options] [path...]")
	// ----------------------------------------------------------
	// Print a List of Supported Options and Their Descriptions
	// ----------------------------------------------------------
	fmt.Println("Options:")
	fmt.Println("  -l   Use long listing format")
	fmt.Println("  -R   List subdirectories recursively")
	fmt.Println("  -a   Include directory entries whose names begin with a dot (.)")
	fmt.Println("  -t   Sort by modification time, newest first")
	fmt.Println("  -r   Reverse order while sorting")
	fmt.Println("  -c   Capture output to file")
	fmt.Println("  -h   Display this help and exit")
}
