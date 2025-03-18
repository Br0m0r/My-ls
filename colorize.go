package main

import (
	"io/fs"
	"os"
)

// ColorizeName returns the file name string with ANSI color codes based on file type.
// It applies specific colors for symlinks, directories, devices, sockets, executables, etc.
func ColorizeName(file fs.DirEntry, info os.FileInfo) string {
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
