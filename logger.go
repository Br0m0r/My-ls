package main

import (
	"log"
	"os"
)

// Initialize sets up the logger.
func Initialize() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// LogError logs an error message with context.
func LogError(err error, context string) {
	if err != nil {
		log.Printf("Error: %s: %v", context, err)
	}
}
