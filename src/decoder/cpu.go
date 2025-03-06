package main

import "fmt"

type CPU struct {
	Registers map[string]int
	Memory    []byte
	Flags     struct {
		Zero bool
		Sign bool
	}
}

var ops = map[Operation]OperationInfo{
	MOV: {
		Execute:      func(_, v int) int { return v },
		WritesResult: true,
		IsArithmetic: false,
	},
	ADD: {
		Execute:      func(c, v int) int { return c + v },
		WritesResult: true,
		IsArithmetic: true,
	},
	SUB: {
		Execute:      func(c, v int) int { return c - v },
		WritesResult: true,
		IsArithmetic: true,
	},
	CMP: {
		Execute:      func(c, v int) int { return c - v },
		WritesResult: false,
		IsArithmetic: true,
	},
}

const sign_mask = 0x8000

func (cpu *CPU) ExecuteReg(operation Operation, dest Register, value int) int {
	if opInfo, ok := ops[operation]; ok {
		current := cpu.Registers[dest.Name]
		newVal := opInfo.Execute(cpu.Registers[dest.Name], value)
		cpu.updateFlags(newVal, opInfo.IsArithmetic)
		if opInfo.WritesResult {
			cpu.Registers[dest.Name] = newVal
			fmt.Printf(" ;  %s:%#x->%#x", dest.Name, current, newVal)
		} else {
			fmt.Printf(" ; ")
		}
		cpu.printFlags()
		return newVal
	}

	panic("Unknown error")
}

func (c *CPU) updateFlags(value int, isArithmetic bool) {
	if !isArithmetic {
		return
	}
	c.Flags.Zero = value == 0
	c.Flags.Sign = (value & sign_mask) != 0
}

func NewCPU(memSize int) *CPU {
	return &CPU{
		Registers: map[string]int{
			"ax": 0, "bx": 0, "cx": 0, "dx": 0,
			"si": 0, "di": 0, "sp": 0, "bp": 0,
		},
		Memory: make([]byte, memSize),
		Flags: struct {
			Zero bool
			Sign bool
		}{
			Zero: false,
			Sign: false,
		},
	}
}
