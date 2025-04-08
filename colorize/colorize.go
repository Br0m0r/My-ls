package colorize

import (
	"io/fs"         // Provides the fs.DirEntry interface.
	"os"            // Used for file mode constants.
	"path/filepath" // Provides functions for manipulating file paths.
	"strings"       // Used for string operations (e.g., ToLower).
)

// ----------------------------------------------------------
// Function: isImageFile
// Purpose: Check if a given file name corresponds to an image file.
// ----------------------------------------------------------
func isImageFile(name string) bool {
	// ----------------------------------------------------------
	// Get File Extension in Lowercase
	// ----------------------------------------------------------
	ext := strings.ToLower(filepath.Ext(name))

	// ----------------------------------------------------------
	// Switch on the file extension to check against known image formats.
	// ----------------------------------------------------------
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}

// ----------------------------------------------------------
// Function: ColorizeName
// Purpose: Return the file name wrapped with ANSI color codes based on file type.
// Parameters:
//
//	file    - the file entry (implements fs.DirEntry)
//	info    - file info (implements os.FileInfo)
//	capture - if true, skip colorization (used when output is captured)
//
// ----------------------------------------------------------
func ColorizeName(file fs.DirEntry, info os.FileInfo, capture bool) string {
	// ----------------------------------------------------------
	// If capturing output, return the plain file name (without color codes)
	// ----------------------------------------------------------
	if capture {
		return file.Name()
	}

	// ----------------------------------------------------------
	// Define ANSI reset code to revert color changes after the file name.
	// ----------------------------------------------------------
	reset := "\033[0m"

	// ----------------------------------------------------------
	// Get the plain file name and its mode (permissions and type)
	// ----------------------------------------------------------
	name := file.Name()
	mode := info.Mode()

	// ----------------------------------------------------------
	// Handle Symbolic Links: Use different colors based on specific names.
	// ----------------------------------------------------------
	if mode&os.ModeSymlink != 0 {
		if name == "fd" {
			// Blue color for "fd"
			return "\033[34m" + name + reset
		}
		if name == "log" {
			// Magenta color for "log"
			return "\033[35m" + name + reset
		}
		// Cyan for other symlinks.
		return "\033[36m" + name + reset
	}

	// ----------------------------------------------------------
	// Handle Directories: Special colors for directories.
	// ----------------------------------------------------------
	if info.IsDir() {
		// For specific directories like "mqueue" or "shm", use a unique style.
		if name == "mqueue" || name == "shm" {
			return "\033[30;42m" + name + reset // Black text on green background.
		}
		// Blue for regular directories.
		return "\033[34m" + name + reset
	}

	// ----------------------------------------------------------
	// Handle Device Files or Named Pipes: Use yellow color.
	// ----------------------------------------------------------
	if mode&(os.ModeDevice|os.ModeNamedPipe) != 0 {
		return "\033[33m" + name + reset
	}

	// ----------------------------------------------------------
	// Handle Sockets: Use magenta color.
	// ----------------------------------------------------------
	if mode&os.ModeSocket != 0 {
		return "\033[35m" + name + reset
	}

	// ----------------------------------------------------------
	// Handle Executable Files: Use green color.
	// ----------------------------------------------------------
	if mode&0111 != 0 {
		return "\033[32m" + name + reset
	}

	// ----------------------------------------------------------
	// Handle Regular Image Files: Use bright magenta (purple) color.
	// ----------------------------------------------------------
	if mode.IsRegular() && isImageFile(name) {
		return "\033[95m" + name + reset
	}

	// ----------------------------------------------------------
	// Default: Return the plain file name without any color modification.
	// ----------------------------------------------------------
	return name
}
