package main

import (
	"eles/ls"
	"os"
)

// main delegates directly to Run.
func main() {
	ls.Run(os.Args[1:])
}
