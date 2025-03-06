package main

import (
	"fmt"
	"strings"
)

func (cpu *CPU) PrintCPU() {
	fmt.Println()
	fmt.Println("Final registers:")
	regs := []string{"ax", "bx", "cx", "dx", "sp", "bp", "si", "di"}
	for _, reg := range regs {
		val := cpu.Registers[reg]
		fmt.Printf("      %s: %#04x (%d)\n", reg, val, val)
	}

	cpu.PrintFlags()
}

func (c *CPU) PrintFlags() {
	flags := ""
	if c.Flags.Sign {
		flags += "S"
	}
	if c.Flags.Zero {
		flags += "Z"
	}

	if flags != "" {
		fmt.Printf("  flags: %s", flags)
	}
}

func (instr *Instruction) PrintInstruction() {
	var b strings.Builder

	b.WriteString(string(instr.Operation))
	b.WriteRune(' ')

	b.WriteString(instr.Dest.String())

	if instr.Src != nil {
		b.WriteRune(',')
		b.WriteRune(' ')
		b.WriteString(instr.Src.String())
	}

	fmt.Print(b.String())
}
