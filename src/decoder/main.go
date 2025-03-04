package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	simulate(os.Args[1])
}
