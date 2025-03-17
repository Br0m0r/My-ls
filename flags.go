package main

import (
	"fmt"
	"os"
)

// Options holds parsed flag values.
type Options struct {
	Long      bool
	Recursive bool
	ShowAll   bool
	TimeSort  bool
	Reverse   bool
	Capture   bool
	Paths     []string
}

// ParseArgs processes command-line arguments.
func ParseArgs(args []string) Options {
	var opts Options
	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
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
		} else {
			opts.Paths = append(opts.Paths, arg)
		}
	}
	// Default to current directory if no paths provided.
	if len(opts.Paths) == 0 {
		opts.Paths = append(opts.Paths, ".")
	}
	return opts
}

// ToMap converts Options to a map[string]bool for backward compatibility.
func (o Options) ToMap() map[string]bool {
	return map[string]bool{
		"l": o.Long,
		"R": o.Recursive,
		"a": o.ShowAll,
		"t": o.TimeSort,
		"r": o.Reverse,
	}
}

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
