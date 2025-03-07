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

	ip := cpu.Registers["ip"]
	for ip < int(len(bytes)) {
		end := min(ip+6, len(bytes))
		instr := ParseInstruction(bytes[ip:end])

		instr.PrintInstruction()

		executor.Execute(instr)
		ip = cpu.Registers["ip"]
		fmt.Println()
	}

	cpu.PrintCPU()
}
