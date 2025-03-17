package main

import (
	"os"
)

// main delegates directly to Run.
func main() {
	Run(os.Args[1:])
}
