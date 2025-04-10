package display

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"eles/colorize"
	"eles/utils"
)

// DisplayFiles prints file names either in long format (if "-l" flag is set)
// or in a compact single-line format.
func DisplayFiles(dirEntries []fs.DirEntry, directoryPath string, optionFlags map[string]bool, outputWriter io.Writer, captureOutput bool) {
	if optionFlags["l"] {
		// For directory listings, always print the "total" line.
		DisplayLongFormat(dirEntries, directoryPath, outputWriter, captureOutput, true)
	} else {
		for index, entry := range dirEntries {
			if index > 0 {
				fmt.Fprint(outputWriter, "  ")
			}
			entryInfo, err := entry.Info()
			if err != nil {
				fmt.Fprint(outputWriter, entry.Name())
			} else {
				fmt.Fprint(outputWriter, colorize.ColorizeName(entry, entryInfo, captureOutput))
			}
		}
		if len(dirEntries) > 0 {
			fmt.Fprintln(outputWriter)
		}
	}
}

// DisplayLongFormat prints detailed file information in a long listing format,
// similar to "ls -l", showing permissions, links, owner, group, size, modification time,
// and file name. The parameter printTotal indicates whether to print the "total" line.
func DisplayLongFormat(dirEntries []fs.DirEntry, directoryPath string, outputWriter io.Writer, captureOutput bool, printTotal bool) {
	maxLinksWidth := 0
	maxOwnerWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0

	// Calculate maximum widths for formatting.
	for _, entry := range dirEntries {
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}
		stat := entryInfo.Sys().(*syscall.Stat_t)
		linksStr := fmt.Sprintf("%d", stat.Nlink)
		if len(linksStr) > maxLinksWidth {
			maxLinksWidth = len(linksStr)
		}
		ownerName := utils.GetOwner(entryInfo)
		if len(ownerName) > maxOwnerWidth {
			maxOwnerWidth = len(ownerName)
		}
		groupName := utils.GetGroup(entryInfo)
		if len(groupName) > maxGroupWidth {
			maxGroupWidth = len(groupName)
		}
		var sizeField string
		if entryInfo.Mode()&os.ModeDevice != 0 {
			major := (stat.Rdev >> 8) & 0xff
			minor := stat.Rdev & 0xff
			sizeField = fmt.Sprintf("%3d, %3d", major, minor)
		} else {
			sizeField = fmt.Sprintf("%d", entryInfo.Size())
		}
		if len(sizeField) > maxSizeWidth {
			maxSizeWidth = len(sizeField)
		}
	}

	// Sum total blocks.
	var totalBlocks int64 = 0
	for _, entry := range dirEntries {
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}
		if stat, ok := entryInfo.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}

	if printTotal {
		fmt.Fprintf(outputWriter, "total %d\n", totalBlocks/2)
	}

	// Print each entry.
	for _, entry := range dirEntries {
		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}
		stat := entryInfo.Sys().(*syscall.Stat_t)
		coloredName := colorize.ColorizeName(entry, entryInfo, captureOutput)
		// If the file is a symlink, append the link target.
		if entry.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(directoryPath, entry.Name()))
			if err == nil {
				coloredName = coloredName + " -> " + linkTarget
			}
		}
		var sizeField string
		if entryInfo.Mode()&os.ModeDevice != 0 {
			major := (stat.Rdev >> 8) & 0xff
			minor := stat.Rdev & 0xff
			sizeField = fmt.Sprintf("%3d, %3d", major, minor)
		} else {
			sizeField = fmt.Sprintf("%d", entryInfo.Size())
		}
		fmt.Fprintf(outputWriter, "%s %*d %-*s %-*s %*s %12s %s\n",
			utils.GetPermissions(entryInfo),
			maxLinksWidth, stat.Nlink,
			maxOwnerWidth, utils.GetOwner(entryInfo),
			maxGroupWidth, utils.GetGroup(entryInfo),
			maxSizeWidth, sizeField,
			formatModTime(entryInfo),
			coloredName)
	}
}

func formatModTime(entryInfo os.FileInfo) string {
	modTime := entryInfo.ModTime()
	now := time.Now()
	sixMonths := time.Hour * 24 * 365 / 2
	if now.Sub(modTime) > sixMonths || modTime.Sub(now) > sixMonths {
		return modTime.Format("Jan _2 2006")
	}
	return modTime.Format("Jan _2 15:04")
}
