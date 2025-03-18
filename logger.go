package main

import (
	"log"
	"os"
)

// Initialize sets up the logger to write log messages to standard error (stderr)
// and configures the log flags to include the date/time and the source file information.
func Initialize() {
	// Set log output destination to standard error.
	log.SetOutput(os.Stderr)
	// Configure log flags to include the date/time (LstdFlags) and the source file and line number (Lshortfile).
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// LogError logs an error message along with some context.
// If err is non-nil, it prints the provided context and the error message.
func LogError(err error, context string) {
	if err != nil {
		// Print an error message with context and the error details.
		log.Printf("Error: %s: %v", context, err)
	}
}
