package main

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
