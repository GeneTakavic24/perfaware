package main

import "fmt"

type CPU struct {
	Registers map[string]int
	Memory    []byte
}

func NewCPU(memSize int) *CPU {
	return &CPU{
		Registers: map[string]int{
			"ax": 0, "bx": 0, "cx": 0, "dx": 0,
			"si": 0, "di": 0, "sp": 0, "bp": 0,
		},
		Memory: make([]byte, memSize),
	}
}

func (cpu *CPU) PrintCPU() {
	fmt.Println()
	fmt.Println("Final registers:")
	regs := []string{"ax", "bx", "cx", "dx", "sp", "bp", "si", "di"} // Order matters
	for _, reg := range regs {
		val := cpu.Registers[reg]
		fmt.Printf("      %s: %#04x (%d)\n", reg, val, val) // %#04x pads to 4 hex digits
	}
}
