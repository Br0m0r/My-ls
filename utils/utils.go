package utils

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"syscall"
)

// PseudoDirEntry implements fs.DirEntry for a single file.
type PseudoDirEntry struct {
	NameStr string
	InfoVal os.FileInfo
}

func (p *PseudoDirEntry) Name() string               { return p.NameStr }
func (p *PseudoDirEntry) IsDir() bool                { return p.InfoVal.IsDir() }
func (p *PseudoDirEntry) Type() fs.FileMode          { return p.InfoVal.Mode().Type() }
func (p *PseudoDirEntry) Info() (fs.FileInfo, error) { return p.InfoVal, nil }

// NewPseudoDirEntry creates a pseudo directory entry from file info.
func NewPseudoDirEntry(info os.FileInfo, name string) fs.DirEntry {
	return &PseudoDirEntry{NameStr: name, InfoVal: info}
}

// GetPermissions returns a string representation of file permissions.
func GetPermissions(info os.FileInfo) string {
	mode := info.Mode()
	var perms string

	// File type.
	switch {
	case mode&os.ModeSymlink != 0:
		perms = "l"
	case mode.IsDir():
		perms = "d"
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			perms = "c"
		} else {
			perms = "b"
		}
	default:
		perms = "-"
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
	// For execute, check setuid.
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
	// For execute, check setgid.
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

	// Others permissions.
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
	// For execute, check sticky bit.
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

// GetOwner returns the owner username of the file.
func GetOwner(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	usr, _ := user.LookupId(fmt.Sprint(uid))
	return usr.Username
}

// GetGroup returns the group name of the file.
func GetGroup(info os.FileInfo) string {
	stat := info.Sys().(*syscall.Stat_t)
	gid := stat.Gid
	grp, _ := user.LookupGroupId(fmt.Sprint(gid))
	return grp.Name
}
