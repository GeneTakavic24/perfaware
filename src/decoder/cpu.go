package main

type CPU struct {
	Registers map[string]int
	Memory    []byte
}

func NewCPU(memSize int) *CPU {
	return &CPU{
		Registers: map[string]int{
			"AX": 0, "BX": 0, "CX": 0, "DX": 0,
			"SI": 0, "DI": 0, "SP": 0, "BP": 0,
		},
		Memory: make([]byte, memSize),
	}
}
