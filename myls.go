// myls.go all in one - first draft
package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Config holds the flag settings for myls.
type Config struct {
	long      bool // -l flag
	recursive bool // -R flag
	all       bool // -a flag
	reverse   bool // -r flag
	sortTime  bool // -t flag
}

// FileEntry holds a file's name and its os.FileInfo.
type FileEntry struct {
	name string
	info os.FileInfo
	path string // full path used when recursing
}

// parseArgs manually parses the command-line arguments.
// It collects flags (starting with '-') and file/directory paths.
func parseArgs(args []string) (Config, []string, error) {
	config := Config{}
	paths := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			// process each flag character (e.g. -lR becomes 'l' and 'R')
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					config.long = true
				case 'R':
					config.recursive = true
				case 'a':
					config.all = true
				case 'r':
					config.reverse = true
				case 't':
					config.sortTime = true
				default:
					return config, nil, errors.New("unsupported flag: -" + string(ch))
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}

	// default to current directory if no path is provided
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	return config, paths, nil
}

// joinPath safely joins two path segments.
func joinPath(dir, file string) string {
	if dir == "/" {
		return "/" + file
	}
	if strings.HasSuffix(dir, "/") {
		return dir + file
	}
	return dir + "/" + file
}

// formatFileMode returns a string representing the file mode similar to ls -l.
func formatFileMode(mode os.FileMode) string {
	var b strings.Builder

	// File type
	switch {
	case mode.IsDir():
		b.WriteByte('d')
	case mode&os.ModeSymlink != 0:
		b.WriteByte('l')
	default:
		b.WriteByte('-')
	}

	// Owner permissions
	if mode&0400 != 0 {
		b.WriteByte('r')
	} else {
		b.WriteByte('-')
	}
	if mode&0200 != 0 {
		b.WriteByte('w')
	} else {
		b.WriteByte('-')
	}
	if mode&0100 != 0 {
		b.WriteByte('x')
	} else {
		b.WriteByte('-')
	}

	// Group permissions
	if mode&0040 != 0 {
		b.WriteByte('r')
	} else {
		b.WriteByte('-')
	}
	if mode&0020 != 0 {
		b.WriteByte('w')
	} else {
		b.WriteByte('-')
	}
	if mode&0010 != 0 {
		b.WriteByte('x')
	} else {
		b.WriteByte('-')
	}

	// Other permissions
	if mode&0004 != 0 {
		b.WriteByte('r')
	} else {
		b.WriteByte('-')
	}
	if mode&0002 != 0 {
		b.WriteByte('w')
	} else {
		b.WriteByte('-')
	}
	if mode&0001 != 0 {
		b.WriteByte('x')
	} else {
		b.WriteByte('-')
	}

	return b.String()
}

// getOwnerGroup returns the user and group names from file stat info.
// If lookup fails, it returns the numeric IDs.
func getOwnerGroup(info os.FileInfo) (string, string) {
	sysStat, ok := info.Sys().(*syscall.Stat_t)
	uid := strconv.Itoa(int(sysStat.Uid))
	gid := strconv.Itoa(int(sysStat.Gid))
	if !ok {
		return uid, gid
	}
	usr, err := user.LookupId(strconv.Itoa(int(sysStat.Uid)))
	grp, err2 := user.LookupGroupId(strconv.Itoa(int(sysStat.Gid)))
	owner := uid
	group := gid
	if err == nil {
		owner = usr.Username
	}
	if err2 == nil {
		group = grp.Name
	}
	return owner, group
}

// formatModTime formats the modification time similar to ls.
// If the file is older than 6 months, it prints the year instead of the time.
func formatModTime(t time.Time) string {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	if t.Before(sixMonthsAgo) || t.After(time.Now()) {
		return t.Format("Jan _2  2006")
	}
	return t.Format("Jan _2 15:04")
}

// sortEntries sorts the FileEntry slice according to the flags.
// It uses a simple insertion sort.
func sortEntries(entries []FileEntry, cfg Config) {
	n := len(entries)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			// decide comparison based on sortTime or lexicographic
			shouldSwap := false
			if cfg.sortTime {
				t1 := entries[j-1].info.ModTime()
				t2 := entries[j].info.ModTime()
				// by default, ls -t sorts newest first
				if t1.Before(t2) {
					shouldSwap = true
				}
			} else {
				if entries[j-1].name > entries[j].name {
					shouldSwap = true
				}
			}
			// if reverse flag is set, invert the swap decision
			if cfg.reverse {
				shouldSwap = !shouldSwap
			}
			if !shouldSwap {
				break
			}
			entries[j-1], entries[j] = entries[j], entries[j-1]
			j--
		}
	}
}

// printLongListing prints a list of FileEntry in long format.
// It computes column widths for link count and file size.
func printLongListing(entries []FileEntry) {
	// Determine maximum widths.
	maxLinks := 0
	maxSize := 0
	linkCounts := make([]uint64, len(entries))
	sizes := make([]int64, len(entries))

	for i, entry := range entries {
		info := entry.info
		// Get link count from syscall.Stat_t
		var links uint64 = 1
		sysStat, ok := info.Sys().(*syscall.Stat_t)
		if ok {
			links = uint64(sysStat.Nlink)
		}
		linkCounts[i] = links
		sizes[i] = info.Size()
		linksLen := len(strconv.FormatUint(links, 10))
		sizeLen := len(strconv.FormatInt(info.Size(), 10))
		if linksLen > maxLinks {
			maxLinks = linksLen
		}
		if sizeLen > maxSize {
			maxSize = sizeLen
		}
	}

	// Print each entry.
	for _, entry := range entries {
		info := entry.info
		modeStr := formatFileMode(info.Mode())
		// Get link count
		var links uint64 = 1
		sysStat, ok := info.Sys().(*syscall.Stat_t)
		if ok {
			links = uint64(sysStat.Nlink)
		}
		owner, group := getOwnerGroup(info)
		modTime := formatModTime(info.ModTime())
		// Format: permissions, links, owner, group, size, mod time, name
		fmt.Printf("%s %*d %-8s %-8s %*d %s %s\n",
			modeStr,
			maxLinks, links,
			owner,
			group,
			maxSize, info.Size(),
			modTime,
			entry.name)
	}
}

// printNames prints just the file names (non-long listing).
func printNames(entries []FileEntry) {
	for _, entry := range entries {
		fmt.Print(entry.name, "  ")
	}
	fmt.Println()
}

// listDirectory lists the contents of a directory given by path.
func listDirectory(path string, cfg Config, printHeader bool) {
	// Open directory.
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "myls: cannot access %s: %v\n", path, err)
		return
	}
	defer f.Close()

	// Read directory entries.
	rawEntries, err := f.Readdir(-1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "myls: error reading directory %s: %v\n", path, err)
		return
	}

	// Build our list of FileEntry.
	var entries []FileEntry
	for _, info := range rawEntries {
		name := info.Name()
		// Skip hidden files if -a not specified.
		if !cfg.all && strings.HasPrefix(name, ".") {
			continue
		}
		entries = append(entries, FileEntry{
			name: name,
			info: info,
			path: joinPath(path, name),
		})
	}

	// Sort entries.
	sortEntries(entries, cfg)

	// If more than one directory is listed or if recursive listing is on, print header.
	if printHeader {
		fmt.Printf("%s:\n", path)
	}

	// Print the entries.
	if cfg.long {
		printLongListing(entries)
	} else {
		printNames(entries)
	}

	// If recursive flag is set, process subdirectories.
	if cfg.recursive {
		for _, entry := range entries {
			if entry.info.IsDir() && entry.name != "." && entry.name != ".." {
				fmt.Println()
				listDirectory(entry.path, cfg, true)
			}
		}
	}
}

// processPath handles each command-line argument: if it is a file, print its info;
// if it is a directory, list its contents.
func processPath(path string, cfg Config, multiple bool) {
	info, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "myls: cannot access %s: %v\n", path, err)
		return
	}

	if !info.IsDir() {
		// For a single file, we want to print the file info.
		fe := FileEntry{name: path, info: info, path: path}
		if cfg.long {
			printLongListing([]FileEntry{fe})
		} else {
			fmt.Println(path)
		}
	} else {
		listDirectory(path, cfg, multiple)
	}
}

func main() {
	cfg, paths, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	multiple := len(paths) > 1
	for i, path := range paths {
		// If more than one path is provided, print header for directories.
		if multiple {
			fmt.Printf("%s:\n", path)
		}
		processPath(path, cfg, multiple)
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
