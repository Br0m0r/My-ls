package main

import (
	"io"
	"log"
	"os"
)

// NewOutput returns an io.Writer that writes to stdout or to a capture file.
func NewOutput(capture bool) io.Writer {
	if capture {
		f, err := os.Create("myls_output.txt")
		if err != nil {
			log.Fatalf("Error creating capture file: %v", err)
		}
		// Write to both stdout and the file.
		return io.MultiWriter(os.Stdout, f)
	}
	return os.Stdout
}
