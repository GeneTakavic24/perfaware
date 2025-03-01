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

	for i := 0; i < len(bytes); {
		// Take up to 6 bytes or whatever's left, decode decides how many to eat
		end := min(i+6, len(bytes))
		instr, consumed := decode(bytes[i:end])
		fmt.Println(instr)
		i += int(consumed) // Move by however many bytes decode used
	}
}
