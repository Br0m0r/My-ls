package main

import (
	"io"
	"log"
	"os"
)

// NewOutput returns an io.Writer that writes to stdout or to a capture file.
// It now returns a cleanup function that should be called after output is done.
func NewOutput(capture bool) (io.Writer, func()) {
	if capture {
		f, err := os.Create("output.txt")
		if err != nil {
			log.Fatalf("Error creating capture file: %v", err)
		}
		// Write to both stdout and the file.
		writer := io.MultiWriter(os.Stdout, f)
		// Return a function that closes the file.
		return writer, func() { f.Close() }
	}
	// For non-capture mode, the cleanup function does nothing.
	return os.Stdout, func() {}
}
