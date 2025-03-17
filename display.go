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
func displayFiles(files []fs.DirEntry, dir string, flags map[string]bool, out io.Writer) {
	if flags["l"] {
		displayLongFormat(files, dir, out)
	} else {
		for i, file := range files {
			if i > 0 {
				fmt.Fprint(out, "  ")
			}
			info, err := file.Info()
			if err != nil {
				fmt.Fprint(out, file.Name())
			} else {
				fmt.Fprint(out, colorizeName(file, info))
			}
		}
		fmt.Fprintln(out)
	}
}

// displayLongFormat prints details in a ls -l style format.
func displayLongFormat(files []fs.DirEntry, dir string, out io.Writer) {
	maxLinksWidth := 0
	maxOwnerWidth := 0
	maxGroupWidth := 0
	maxSizeWidth := 0

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
	fmt.Fprintf(out, "total %d\n", totalBlocks/2)

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		stat := info.Sys().(*syscall.Stat_t)
		coloredName := colorizeName(file, info)
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
		fmt.Fprintf(out, "%s %*d %-*s %-*s %*s %12s %s\n",
			getPermissions(info),
			maxLinksWidth, stat.Nlink,
			maxOwnerWidth, getOwner(info),
			maxGroupWidth, getGroup(info),
			maxSizeWidth, sizeField,
			formatModTime(info),
			coloredName)
	}
}

func formatModTime(info os.FileInfo) string {
	modTime := info.ModTime()
	now := time.Now()
	sixMonths := time.Hour * 24 * 365 / 2
	if now.Sub(modTime) > sixMonths || modTime.Sub(now) > sixMonths {
		return modTime.Format("Jan _2 2006")
	}
	return modTime.Format(" Jan _2 15:04")
}
