package main

import (
	"fmt"
	"os"
	"strings"
)

// Run handles command-line arguments and calls appropriate functions.
func Run(args []string) {
	dir := "." // Default directory
	flags := parseFlags(args, &dir)

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: No such file or directory:", dir)
		} else if os.IsPermission(err) {
			fmt.Println("Error: Permission denied:", dir)
		} else {
			fmt.Println("Error:", err)
		}
		return
	}

	// If it's a file, print its name (ignoring `-R`)
	if !info.IsDir() {
		fmt.Println(dir)
		return
	}

	// Call recursive function if `-R` flag is enabled
	if flags["R"] {
		recursiveList(dir, flags)
	} else {
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		files = filterFiles(files, flags, dir)

		files = sortFiles(files, flags)
		displayFiles(files, dir, flags)
	}
}

// parseFlags extracts flags and updates the directory argument.
func parseFlags(args []string, dir *string) map[string]bool {
	flags := map[string]bool{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, char := range arg[1:] {
				flags[string(char)] = true
			}
		} else {
			*dir = arg // Directory argument
		}
	}

	return flags
}
