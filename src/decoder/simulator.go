package main

import (
	"fmt"
	"os"
)

func simulate(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	cpu := NewCPU(10)

	executor := newX86Executor(cpu)

	fmt.Printf("; %s disassembly:\n", filePath)
	fmt.Println("bits 16")
	fmt.Println()

	for i := 0; i < len(bytes); {
		end := min(i+6, len(bytes))
		instr := ParseInstruction(bytes[i:end])
		instr.PrintInstruction()
		executor.Execute(instr)
		fmt.Println()
		i += int(instr.Consumed)
	}

	cpu.PrintCPU()
}
