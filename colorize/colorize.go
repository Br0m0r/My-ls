package colorize

import (
	"io/fs"         
	"os"            
	"path/filepath" 
	"strings"       
)

//Check if a given file name corresponds to an image file.
func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}


//Return the file name wrapped with ANSI color codes based on file type.
func ColorizeName(file fs.DirEntry, info os.FileInfo, capture bool) string {
	if capture {
		return file.Name()
	}
	reset := "\033[0m"
	name := file.Name()
	mode := info.Mode() //permissions and type

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

	if info.IsDir() {
		if name == "mqueue" || name == "shm" {
			return "\033[30;42m" + name + reset // For specific directories like "mqueue" or "shm" . Black text on green background.
		}
		// Blue for regular directories.
		return "\033[34m" + name + reset
	}

	if mode&(os.ModeDevice|os.ModeNamedPipe) != 0 {   //Device Files or Named Pipes: Use yellow color
		return "\033[33m" + name + reset
	}

	if mode&os.ModeSocket != 0 {  // Handle Sockets: Use magenta color.
		return "\033[35m" + name + reset
	}
	
	if mode&0111 != 0 {                   //  Executable Files: Use green color
		return "\033[32m" + name + reset
	}

	if mode.IsRegular() && isImageFile(name) {            // Regular Image Files: Use bright magenta (purple) color.
		return "\033[95m" + name + reset
	}
	
	return name     //plain file
}
