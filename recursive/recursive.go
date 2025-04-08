package recursive

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"eles/display"
	"eles/filter"
	"eles/sort"
)

// joinDisplayPath joins parent and child directory names, preserving the "./" prefix when the parent is "." or starts with "./".
func joinDisplayPath(parent, child string) string {
	if parent == "." {
		return "./" + child
	}
	if strings.HasPrefix(parent, "./") {
		if strings.HasSuffix(parent, "/") {
			return parent + child
		}
		return parent + "/" + child
	}
	return filepath.Join(parent, child)
}

// RecursiveList lists directories recursively.
func RecursiveList(dir string, flags map[string]bool, capture bool, out io.Writer) {
	fmt.Fprintf(out, "\n%s:\n", dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(out, "Error: %v\n", err)
		return
	}

	files = filter.FilterFiles(files, flags, dir)
	files = sort.SortFiles(files, flags)
	display.DisplayFiles(files, dir, flags, out, capture)

	for _, file := range files {
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		if file.Type()&os.ModeSymlink != 0 {
			continue
		}
		if info.IsDir() {
			subDir := joinDisplayPath(dir, file.Name())
			RecursiveList(subDir, flags, capture, out)
		}
	}
}
