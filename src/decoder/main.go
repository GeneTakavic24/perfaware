package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	decode(os.Args[1])
}
