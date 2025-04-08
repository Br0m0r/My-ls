package output

import (
	"io"  // Provides interfaces for I/O operations.
	"log" // Used for logging errors.
	"os"  // Provides OS functionality like file creation.
)

// NewOutput returns an io.Writer that writes to stdout or to a capture file,
// along with a cleanup function to close any resources.
func NewOutput(capture bool) (io.Writer, func()) {
	// ----------------------------------------------------------
	// Check if Output Capture is Enabled (-c flag)
	// ----------------------------------------------------------
	if capture {
		// ----------------------------------------------------------
		// Create or Overwrite a File Named "output.txt"
		// ----------------------------------------------------------
		f, err := os.Create("output.txt")
		// If file creation fails, log the error and terminate.
		if err != nil {
			log.Fatalf("Error creating capture file: %v", err)
		}
		// ----------------------------------------------------------
		// Create a MultiWriter: Output Goes to Both Stdout and the Capture File
		// ----------------------------------------------------------
		writer := io.MultiWriter(os.Stdout, f)
		// Return the writer and a cleanup function to close the file later.
		return writer, func() { f.Close() }
	}
	// ----------------------------------------------------------
	// If Capture is Disabled, Use Standard Output Only
	// ----------------------------------------------------------
	return os.Stdout, func() {} // An empty cleanup function since there's nothing to close.
}
