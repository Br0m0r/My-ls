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

// DisplayFiles prints file names either in long format (if flag "l" is set) or in a single line.
func DisplayFiles(files []fs.DirEntry, dir string, flags map[string]bool, out io.Writer, capture bool) {
	if flags["l"] {
		DisplayLongFormat(files, dir, out, capture)
	} else {
		for i, file := range files {
			if i > 0 {
				fmt.Fprint(out, "  ")
			}
			info, err := file.Info()
			if err != nil {
				fmt.Fprint(out, file.Name())
			} else {
				fmt.Fprint(out, colorize.ColorizeName(file, info, capture))
			}
		}
		fmt.Fprintln(out)
	}
}

// DisplayLongFormat prints detailed file information (like "ls -l").
func DisplayLongFormat(files []fs.DirEntry, dir string, out io.Writer, capture bool) {
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

		owner := utils.GetOwner(info)
		if len(owner) > maxOwnerWidth {
			maxOwnerWidth = len(owner)
		}

		group := utils.GetGroup(info)
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
		coloredName := colorize.ColorizeName(file, info, capture)
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
			utils.GetPermissions(info),
			maxLinksWidth, stat.Nlink,
			maxOwnerWidth, utils.GetOwner(info),
			maxGroupWidth, utils.GetGroup(info),
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
