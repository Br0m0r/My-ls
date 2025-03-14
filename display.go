package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// displayFiles prints file names either in long format or in a single horizontal line.
func displayFiles(files []fs.DirEntry, dir string, flags map[string]bool) {
	if flags["l"] {
		displayLongFormat(files, dir)
	} else {
		for i, file := range files {
			if i > 0 {
				fmt.Print("  ")
			}
			info, err := file.Info()
			if err != nil {
				fmt.Print(file.Name())
			} else {
				fmt.Print(colorizeName(file, info))
			}
		}
		fmt.Println()
	}
}

// displayLongFormat prints details in a ls -l style format.
// It first computes the maximum widths of several columns so that the output is aligned.
func displayLongFormat(files []fs.DirEntry, dir string) {
	// First pass: compute maximum widths for links, owner, group, and size fields.
	maxLinksWidth := 0
	maxOwnerWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0

	// We'll store each file's FileInfo in a slice to avoid re-fetching it.
	var fileInfos []os.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, info)

		stat := info.Sys().(*syscall.Stat_t)
		linksStr := fmt.Sprintf("%d", stat.Nlink)
		if len(linksStr) > maxLinksWidth {
			maxLinksWidth = len(linksStr)
		}

		owner := getOwner(info)
		if len(owner) > maxOwnerWidth {
			maxOwnerWidth = len(owner)
		}

		group := getGroup(info)
		if len(group) > maxGroupWidth {
			maxGroupWidth = len(group)
		}

		var sizeField string
		if info.Mode()&os.ModeDevice != 0 {
			// For device files, display major and minor numbers.
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

	// Print total block count (divided by 2 as ls does) only if not a single file.
	if !(len(files) == 1 && !files[0].IsDir()) {
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
		fmt.Printf("total %d\n", totalBlocks/2)
	}

	// Second pass: print each file's details using the computed widths.
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)
		coloredName := colorizeName(file, info)
		// If the file is a symlink, append its target.
		if file.Type()&os.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(filepath.Join(dir, file.Name()))
			if err == nil {
				coloredName = coloredName + " -> " + linkTarget
			}
		}
		var sizeField string
		if info.Mode()&os.ModeDevice != 0 {
			major := (stat.Rdev >> 8) & 0xff
			minor := stat.Rdev & 0xff
			sizeField = fmt.Sprintf("%3d, %3d", major, minor)
		} else {
			sizeField = fmt.Sprintf("%d", info.Size())
		}

		// Print using dynamic width formatting.
		fmt.Printf("%s %*d %-*s %-*s %*s %12s %s\n",
			getPermissions(info),
			maxLinksWidth, stat.Nlink,
			maxOwnerWidth, getOwner(info),
			maxGroupWidth, getGroup(info),
			maxSizeWidth, sizeField,
			formatModTime(info),
			coloredName)
	}
}

// formatModTime formats the modification time like ls -l.
func formatModTime(info os.FileInfo) string {
	modTime := info.ModTime()
	now := time.Now()
	sixMonths := time.Hour * 24 * 365 / 2
	if now.Sub(modTime) > sixMonths || modTime.Sub(now) > sixMonths {
		return modTime.Format("Jan _2 2006")
	}
	return modTime.Format("Jan _2 15:04")
}
