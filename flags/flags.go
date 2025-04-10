package flags

import (
	"fmt"
	"os"
)

// Options holds parsed flag values and file/directory paths.
type Options struct {
	Long      bool     //(-l)
	Recursive bool     //  (-R)
	ShowAll   bool     // (-a)
	TimeSort  bool     //  (-t)
	Reverse   bool     // (-r)
	Capture   bool     // (-c)
	Paths     []string 
}

func ParseArgs(args []string) Options {
	var opts Options
	endOfOptions := false

	for _, arg := range args {
		
		if !endOfOptions && len(arg) > 0 && arg[0] == '-' {
				if arg == "--" {
				endOfOptions = true
				continue
			}
			
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

	// if no paths provided, use current directory.
	if len(opts.Paths) == 0 {
		opts.Paths = append(opts.Paths, ".")
	}
	return opts
}

// ToMap converts Options to a map for compatibility with other functions.
func (o Options) ToMap() map[string]bool {
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
	fmt.Println("Usage: myls [options] [path...]")
	fmt.Println("Options:")
	fmt.Println("  -l   Use long listing format")
	fmt.Println("  -R   List subdirectories recursively")
	fmt.Println("  -a   Include directory entries whose names begin with a dot (.)")
	fmt.Println("  -t   Sort by modification time, newest first")
	fmt.Println("  -r   Reverse order while sorting")
	fmt.Println("  -c   Capture output to file")
	fmt.Println("  -h   Display this help and exit")
}
