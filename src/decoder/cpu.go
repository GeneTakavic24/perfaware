package main

import (
	"fmt"
)

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

func (cpu *CPU) ExecuteJump(operation Operation, offset int) bool {
	switch operation {
	case Jnz:
		if !cpu.Flags.Zero {
			oldIP := cpu.Registers["ip"]
			cpu.Registers["ip"] = oldIP + offset
			return true
		}
	}
	return false
}

func (cpu *CPU) ExecuteReg(operation Operation, dest Register, value int) int {
	if opInfo, ok := ops[operation]; ok {
		current := cpu.Registers[dest.Name]
		newVal := opInfo.Execute(cpu.Registers[dest.Name], value)
		if opInfo.WritesResult {
			cpu.Registers[dest.Name] = newVal
			fmt.Printf(" ;  %s:%#x->%#x", dest.Name, current, newVal)
		} else {
			fmt.Printf(" ; ")
		}

		cpu.updateFlags(newVal, opInfo.IsArithmetic)
		return newVal
	}

	panic("Unknown error")
}

func (cpu *CPU) ExecuteMem(operation Operation, dest EffectiveAddress, value int) int {
	if opInfo, ok := ops[operation]; ok {
		high, low := cpu.resolveEffectiveAddress(dest)
		h := uint(*high) << 8
		combined := int(h | uint(*low))
		newVal := opInfo.Execute(combined, value)

		if opInfo.WritesResult {
			*low = byte(newVal & 0xFF)
			*high = byte((newVal >> 8) & 0xFF)
		}

		cpu.updateFlags(newVal, opInfo.IsArithmetic)
		return newVal
	}

	panic("Unknown error")
}

func (c *CPU) updateFlags(value int, isArithmetic bool) {
	if !isArithmetic {
		return
	}
	oldZero := c.Flags.Zero
	oldSign := c.Flags.Sign

	c.Flags.Zero = value == 0
	c.Flags.Sign = (value & sign_mask) != 0

	if !oldZero && c.Flags.Zero {
		fmt.Printf(" flags:->Z")
	} else if oldZero && !c.Flags.Zero {
		fmt.Printf(" flags:Z->")
	}

	if !oldSign && c.Flags.Sign {
		fmt.Printf(" flags:->S")
	} else if oldSign && !c.Flags.Sign {
		fmt.Printf(" flags:S->")
	}
}

func NewCPU(memSize uint16) *CPU {
	return &CPU{
		Registers: map[string]int{
			"ax": 0, "bx": 0, "cx": 0, "dx": 0,
			"si": 0, "di": 0, "sp": 0, "bp": 0,
			"ip": 0,
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

func (cpu *CPU) resolveEffectiveAddress(addr EffectiveAddress) (high, low *byte) {
	address := int(addr.Offset)
	if addr.Base.Name != "" {
		address += cpu.Registers[addr.Base.Name]
	}
	if addr.Index.Name != "" {
		address += cpu.Registers[addr.Index.Name]
	}

	lowByte := &cpu.Memory[address]
	highByte := &cpu.Memory[address+1]
	return highByte, lowByte
}
