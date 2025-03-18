package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// displayFiles prints file names either in long format or in a single horizontal line.
// The function takes in:
// - files: a slice of directory entries to display.
// - dir: the directory where the files are located.
// - flags: a map of command-line flags that affect how files are displayed.
// - out: the io.Writer where the output should be written.
// - capture: a flag that, if true, disables ANSI color codes (useful for output capture).
func displayFiles(files []fs.DirEntry, dir string, flags map[string]bool, out io.Writer, capture bool) {
	// If the long listing flag (-l) is set, use the long format display.
	if flags["l"] {
		displayLongFormat(files, dir, out, capture)
	} else {
		// Otherwise, display file names on one line separated by two spaces.
		for i, file := range files {
			// Add a space separator between files (but not before the first file).
			if i > 0 {
				fmt.Fprint(out, "  ")
			}
			// Try to get file information; if an error occurs, print the file name directly.
			info, err := file.Info()
			if err != nil {
				fmt.Fprint(out, file.Name())
			} else {
				// ColorizeName applies ANSI color codes based on file type,
				// unless the capture flag is set (in which case no color is applied).
				fmt.Fprint(out, ColorizeName(file, info, capture))
			}
		}
		// End the line after printing all file names.
		fmt.Fprintln(out)
	}
}

// displayLongFormat prints detailed file information in a style similar to "ls -l".
// It calculates maximum widths for columns (links, owner, group, size) to align the output properly.
func displayLongFormat(files []fs.DirEntry, dir string, out io.Writer, capture bool) {
	// Variables for the maximum width of certain fields for proper alignment.
	maxLinksWidth := 0
	maxOwnerWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0

	// Temporary slice to store FileInfo for each file.
	var fileInfos []os.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, info)

		// Retrieve file system specific information.
		stat := info.Sys().(*syscall.Stat_t)
		// Convert the number of links to string to determine its width.
		linksStr := fmt.Sprintf("%d", stat.Nlink)
		if len(linksStr) > maxLinksWidth {
			maxLinksWidth = len(linksStr)
		}

		// Get the owner's username and update the maximum width.
		owner := GetOwner(info)
		if len(owner) > maxOwnerWidth {
			maxOwnerWidth = len(owner)
		}

		// Get the group name and update the maximum width.
		group := GetGroup(info)
		if len(group) > maxGroupWidth {
			maxGroupWidth = len(group)
		}

		// Prepare the size field. For device files, display major and minor numbers.
		var sizeField string
		if info.Mode()&os.ModeDevice != 0 {
			major := (stat.Rdev >> 8) & 0xff
			minor := stat.Rdev & 0xff
			sizeField = fmt.Sprintf("%3d, %3d", major, minor)
		} else {
			sizeField = fmt.Sprintf("%d", info.Size())
		}
		if len(sizeField) > maxSizeWidth {
			maxSizeWidth = len(sizeField)
		}
	}

	// Calculate the total number of blocks used by the files.
	var totalBlocks int64 = 0
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}
	// Print the total blocks (divided by 2 to match the traditional ls output).
	fmt.Fprintf(out, "total %d\n", totalBlocks/2)

	// Loop over each file again to print detailed information.
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)
		// Get the file name with possible ANSI colorization.
		coloredName := ColorizeName(file, info, capture)
		// If the file is a symbolic link, resolve and append its target.
		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(dir, file.Name()))
			if err == nil {
				coloredName = coloredName + " -> " + linkTarget
			}
		}
		// Prepare the size field, handling device files separately.
		var sizeField string
		if info.Mode()&os.ModeDevice != 0 {
			major := (stat.Rdev >> 8) & 0xff
			minor := stat.Rdev & 0xff
			sizeField = fmt.Sprintf("%3d, %3d", major, minor)
		} else {
			sizeField = fmt.Sprintf("%d", info.Size())
		}
		// Print the formatted long listing:
		// Permissions, number of links, owner, group, size, modification time, and file name.
		fmt.Fprintf(out, "%s %*d %-*s %-*s %*s %12s %s\n",
			GetPermissions(info),      // Permissions (e.g., -rw-r--r--)
			maxLinksWidth, stat.Nlink, // Number of links, right-aligned
			maxOwnerWidth, GetOwner(info), // Owner's username, left-aligned
			maxGroupWidth, GetGroup(info), // Group name, left-aligned
			maxSizeWidth, sizeField, // File size (or device numbers), right-aligned
			formatModTime(info), // Formatted modification time
			coloredName)         // File name (with ANSI color codes if applicable)
	}
}

// formatModTime formats the modification time of a file.
// If the file's modification time is older than six months (or in the future by that amount),
// it prints the date with the year; otherwise, it prints the date with the time.
func formatModTime(info os.FileInfo) string {
	modTime := info.ModTime()
	now := time.Now()
	// Define a six-month duration.
	sixMonths := time.Hour * 24 * 365 / 2
	// If the modification time is more than six months away from now, format it with the year.
	if now.Sub(modTime) > sixMonths || modTime.Sub(now) > sixMonths {
		return modTime.Format("Jan _2 2006")
	}
	// Otherwise, include the time.
	return modTime.Format(" Jan _2 15:04")
}
