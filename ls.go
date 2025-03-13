package main

import (
	"fmt"
	"os"
	"strings"
)

// Run handles command-line arguments and calls appropriate functions.
func Run(args []string) {
	flags, paths := parseFlags(args)
	// Default to current directory if no paths are provided.
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	multiple := len(paths) > 1

	for i, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Error: No such file or directory:", path)
			} else if os.IsPermission(err) {
				fmt.Println("Error: Permission denied:", path)
			} else {
				fmt.Println("Error:", err)
			}
			continue
		}

		if multiple {
			fmt.Println(path + ":")
		}

		// If it's a file, print its name (ignoring -R)
		if !info.IsDir() {
			fmt.Println(path)
		} else {
			if flags["R"] {
				recursiveList(path, flags)
			} else {
				files, err := os.ReadDir(path)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				files = filterFiles(files, flags, path)
				files = sortFiles(files, flags)
				displayFiles(files, path, flags)
			}
		}

		// Print an empty line between listings for multiple paths.
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}

// parseFlags extracts flags and collects non-flag arguments as paths.
func parseFlags(args []string) (map[string]bool, []string) {
	flags := map[string]bool{}
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, char := range arg[1:] {
				flags[string(char)] = true
			}
		} else {
			paths = append(paths, arg)
		}
	}
	return flags, paths
}
