package main

import (
	"log"
	"os"
)

// Initialize sets up the logger to write to stderr with date/time and source file info.
func Initialize() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// LogError logs an error message with context if err is non-nil.
func LogError(err error, context string) {
	if err != nil {
		log.Printf("Error: %s: %v", context, err)
	}
}
