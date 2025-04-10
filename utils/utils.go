package utils

import (
	"fmt"    
	"io/fs"   
	"os"      
	"os/user"
	"syscall" 
)


//Implements the fs.DirEntry interface for a file,
type PseudoDirEntry struct {
	NameStr string      // The name to be displayed.
	InfoVal os.FileInfo // The file information associated with the entry.
}

//Returns the stored name of the pseudo directory entry.
func (p *PseudoDirEntry) Name() string {
	return p.NameStr
}

//Returns true if the underlying file info indicates a directory.
func (p *PseudoDirEntry) IsDir() bool {
	return p.InfoVal.IsDir()
}

//Returns the file mode type extracted from the underlying file info.
func (p *PseudoDirEntry) Type() fs.FileMode {
	return p.InfoVal.Mode().Type()
}

// Returns the underlying os.FileInfo and nil error, fulfilling fs.DirEntry.
func (p *PseudoDirEntry) Info() (fs.FileInfo, error) {
	return p.InfoVal, nil
}


// Creates and returns a pseudo directory entry from given file info and a name.
func NewPseudoDirEntry(info os.FileInfo, name string) fs.DirEntry {
	return &PseudoDirEntry{
		NameStr: name,
		InfoVal: info,
	}
}


// Generates a string representing the file's permissions in a format
func GetPermissions(info os.FileInfo) string {
	mode := info.Mode()
	var perms string

	
	// Determine the file type indicator.
	switch {
	case mode&os.ModeSymlink != 0:
		perms = "l" // Symbolic link.
	case mode.IsDir():
		perms = "d" // Directory.
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			perms = "c" // Character device.
		} else {
			perms = "b" // Block device.
		}
	default:
		perms = "-" // Regular file.
	}

	
	// Owner permissions.
	if mode&0400 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0200 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}

	// Check execute permission and setuid.
	if mode&os.ModeSetuid != 0 {
		if mode&0100 != 0 {
			perms += "s"
		} else {
			perms += "S"
		}
	} else {
		if mode&0100 != 0 {
			perms += "x"
		} else {
			perms += "-"
		}
	}

	// Group permissions.
	if mode&0040 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0020 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	// Check execute permission and setgid.
	if mode&os.ModeSetgid != 0 {
		if mode&0010 != 0 {
			perms += "s"
		} else {
			perms += "S"
		}
	} else {
		if mode&0010 != 0 {
			perms += "x"
		} else {
			perms += "-"
		}
	}

	if mode&0004 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0002 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	// Check execute permission and sticky bit.
	if mode&os.ModeSticky != 0 {
		if mode&0001 != 0 {
			perms += "t"
		} else {
			perms += "T"
		}
	} else {
		if mode&0001 != 0 {
			perms += "x"
		} else {
			perms += "-"
		}
	}

	return perms
}

// Retrieves and returns the owner username for the file based on its UID.
func GetOwner(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)     // Convert system-specific data to *syscall.Stat_t.
	uid := stat.Uid                          // Retrieve the UID from the stat.
	usr, _ := user.LookupId(fmt.Sprint(uid)) // Look up the username using the UID.
	return usr.Username                      // Return the username.
}


// Retrieves and returns the group name for the file based on its GID.
func GetGroup(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)          // Convert system-specific data to *syscall.Stat_t.
	gid := stat.Gid                               // Retrieve the GID from the stat.
	grp, _ := user.LookupGroupId(fmt.Sprint(gid)) // Look up the group name using the GID.
	return grp.Name                               // Return the group name.
}
