package output

import (
	"io"  
	"log" 
	"os"  
)

// NewOutput returns an io.Writer that writes to stdout or to a capture file,
// along with a cleanup function to close any resources.
func NewOutput(capture bool) (io.Writer, func()) {
	if capture {
		f, err := os.Create("output.txt")
		if err != nil {
		log.Fatalf("Error creating capture file: %v", err)
		}
		writer := io.MultiWriter(os.Stdout, f)
			return writer, func() { f.Close() }
	}
	
	return os.Stdout, func() {} // An empty cleanup function since there's nothing to close.
}
