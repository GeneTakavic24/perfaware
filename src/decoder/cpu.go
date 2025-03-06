package main

type CPU struct {
	Registers map[string]int
	Memory    []byte
	Flags     struct {
		Zero bool
		Sign bool
	}
}

const sign_mask = 0x8000

func (cpu *CPU) Execute(opInfo OperationInfo, regName string, val int) int {
	newVal := opInfo.Execute(cpu.Registers[regName], val)
	cpu.updateFlags(newVal, opInfo.IsArithmetic)
	if opInfo.WritesResult {
		cpu.Registers[regName] = newVal
	}
	return newVal
}

func (c *CPU) updateFlags(value int, isArithmetic bool) {
	if !isArithmetic {
		return
	}
	c.Flags.Zero = value == 0
	c.Flags.Sign = value&sign_mask == 1
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
