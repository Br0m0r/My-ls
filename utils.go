package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"syscall"
)

// getPermissions returns a string representation of file permissions.
func getPermissions(info os.FileInfo) string {
	mode := info.Mode()
	perms := ""

	if mode.IsDir() {
		perms += "d"
	} else {
		perms += "-"
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

// getOwner returns the owner of the file.
func getOwner(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	user, _ := user.LookupId(fmt.Sprint(uid))
	return user.Username
}

// getGroup returns the group of the file.
func getGroup(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	gid := stat.Gid
	group, _ := user.LookupGroupId(fmt.Sprint(gid))
	return group.Name
}

// colorizeName returns the file name string with ANSI color codes based on the file type.
func colorizeName(file fs.DirEntry, info os.FileInfo) string {
	// Default color reset code
	reset := "\033[0m"

	// Check for symbolic links first.
	if file.Type()&os.ModeSymlink != 0 {
		return "\033[36m" + file.Name() + reset
	}

	// Directory: Blue
	if info.IsDir() {
		return "\033[34m" + file.Name() + reset
	}

	// Check for executables (if any execute bits are set)
	mode := info.Mode()
	if mode&0111 != 0 { // any execute bit is set
		return "\033[32m" + file.Name() + reset
	}

	// Default for regular files
	return file.Name()
}
