package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"syscall"
)

// getPermissions returns a string representation of file permissions,
// using "d" for directories, "c" for character devices, "b" for block devices, and "-" otherwise.
func getPermissions(info os.FileInfo) string {
	mode := info.Mode()
	var perms string
	if mode.IsDir() {
		perms = "d"
	} else if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			perms = "c"
		} else {
			perms = "b"
		}
	} else {
		perms = "-"
	}

	permBits := []rune{'r', 'w', 'x'}
	for i := 0; i < 9; i++ {
		if mode&(1<<uint(8-i)) != 0 {
			perms += string(permBits[i%3])
		} else {
			perms += "-"
		}
	}

	return perms
}

// getOwner returns the owner username of the file.
func getOwner(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	usr, _ := user.LookupId(fmt.Sprint(uid))
	return usr.Username
}

// getGroup returns the group name of the file.
func getGroup(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	gid := stat.Gid
	grp, _ := user.LookupGroupId(fmt.Sprint(gid))
	return grp.Name
}

// colorizeName returns the file name string with ANSI color codes based on file type.
// The following adjustments have been made for /dev:
//   - Symlink "fd" is shown in blue.
//   - Symlink "log" is shown in magenta.
//   - Directories "mqueue" and "shm" are shown with a green background (black text).
//
// Other rules remain: symlinks (default) are cyan, directories are blue,
// device files/named pipes are yellow, sockets are magenta, and executables are green.
func colorizeName(file fs.DirEntry, info os.FileInfo) string {
	reset := "\033[0m"
	name := file.Name()
	mode := info.Mode()

	// For symlinks: check special cases for "fd" and "log"
	if mode&os.ModeSymlink != 0 {
		if name == "fd" {
			// Blue for the "fd" symlink (to /proc/self/fd)
			return "\033[34m" + name + reset
		}
		if name == "log" {
			// Magenta for the "log" symlink (to /run/systemd/journal/dev-log)
			return "\033[35m" + name + reset
		}
		// Default symlink color: cyan
		return "\033[36m" + name + reset
	}

	// For directories: if the name is "mqueue" or "shm", use a green background.
	if info.IsDir() {
		if name == "mqueue" || name == "shm" {
			// Green background with black text
			return "\033[30;42m" + name + reset
		}
		// Default directory color: blue.
		return "\033[34m" + name + reset
	}

	// Device files or named pipes (FIFO): yellow.
	if mode&(os.ModeDevice|os.ModeNamedPipe) != 0 {
		return "\033[33m" + name + reset
	}

	// Sockets: magenta.
	if mode&os.ModeSocket != 0 {
		return "\033[35m" + name + reset
	}

	// Executables: green.
	if mode&0111 != 0 {
		return "\033[32m" + name + reset
	}

	// Default: no color.
	return name
}
