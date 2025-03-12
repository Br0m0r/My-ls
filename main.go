package main

import (
	"os"
)

func main() {
	Run(os.Args[1:]) // Delegates everything to ls package
}
