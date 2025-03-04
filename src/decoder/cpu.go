package main

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
