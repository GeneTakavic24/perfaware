package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	filePath := os.Args[1]

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("; %s disassembly:\n", filePath)
	fmt.Println("bits 16")
	fmt.Println()

	for i := 0; i < len(bytes); i += 2 {
		pair := bytes[i : i+2]
		fmt.Println(decode(&pair))
	}
}
