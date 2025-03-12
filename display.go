package main

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
	"time"
)

// displayFiles prints file names in a **single horizontal line**.
func displayFiles(files []fs.DirEntry, dir string, flags map[string]bool) {
	if flags["l"] {
		displayLongFormat(files, dir)
	} else {
		for i, file := range files {
			if i > 0 {
				fmt.Print("  ")
			}
			// Get the file info to decide on color.
			info, err := file.Info()
			if err != nil {
				// Fallback to no color on error.
				fmt.Print(file.Name())
			} else {
				fmt.Print(colorizeName(file, info))
			}
		}
		fmt.Println()
	}
}

// displayLongFormat prints details in `ls -l` format.
func displayLongFormat(files []fs.DirEntry, dir string) {
	var totalBlocks int64 = 0
	// Compute the total blocks for the files in the directory.
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}
	// The real ls often divides the block count by 2 (for 1K blocks vs. 512-byte blocks)
	fmt.Printf("total %d\n", totalBlocks/2)

	// Now display each file's details in long format.
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)
		// Use colorizeName for file name output
		coloredName := colorizeName(file, info)
		fmt.Printf("%s %2d %s %s %6d %s %s\n",
			getPermissions(info),
			stat.Nlink,
			getOwner(info),
			getGroup(info),
			info.Size(),
			formatModTime(info),
			coloredName,
		)
	}
}

// formatModTime formats the modification time like `ls -l`.
func formatModTime(info os.FileInfo) string {
	modTime := info.ModTime()
	currentYear := time.Now().Year()

	if modTime.Year() == currentYear {
		return modTime.Format("Jan 02 15:04")
	} else {
		return modTime.Format("Jan 02 2006")
	}
}
