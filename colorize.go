package main

import (
	"io/fs"
	"os"
)

// ColorizeName returns the file name with ANSI color codes based on its type.
func ColorizeName(file fs.DirEntry, info os.FileInfo, capture bool) string {
	if capture {
		return file.Name()
	}
	reset := "\033[0m"
	name := file.Name()
	mode := info.Mode()

	// Handle symbolic links.
	if mode&os.ModeSymlink != 0 {
		if name == "fd" {
			return "\033[34m" + name + reset
		}
		if name == "log" {
			return "\033[35m" + name + reset
		}
		return "\033[36m" + name + reset
	}

	// Handle directories.
	if info.IsDir() {
		if name == "mqueue" || name == "shm" {
			return "\033[30;42m" + name + reset
		}
		return "\033[34m" + name + reset
	}

	// Handle device files or named pipes.
	if mode&(os.ModeDevice|os.ModeNamedPipe) != 0 {
		return "\033[33m" + name + reset
	}

	// Handle sockets.
	if mode&os.ModeSocket != 0 {
		return "\033[35m" + name + reset
	}

	// Handle executables.
	if mode&0111 != 0 {
		return "\033[32m" + name + reset
	}

	return name
}
