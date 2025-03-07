package main

import "fmt"

type Executor interface {
	Execute(instr Instruction) error
}

type X86Executor struct {
	cpu *CPU
}

type OperationInfo struct {
	Execute      func(current, value int) int
	WritesResult bool
	IsArithmetic bool
}

func newX86Executor(cpu *CPU) *X86Executor {
	return &X86Executor{cpu: cpu}
}

func (e *X86Executor) Execute(instr Instruction) error {
	value := e.extractFrom(instr.Src)

	prevIp := e.cpu.Registers["ip"]
	e.cpu.Registers["ip"] = prevIp + int(instr.Consumed)

	defer func() {
		if e.cpu.Registers["ip"] != prevIp {
			fmt.Printf("  ip:0x%x->0x%x", prevIp, e.cpu.Registers["ip"])
		}
	}()

	if dest, ok := instr.Dest.(Register); ok {
		e.cpu.ExecuteReg(instr.Operation, dest, value)
	}

	return nil
}

func (e *X86Executor) extractFrom(o Operand) int {
	switch v := o.(type) {
	case Register:
		return e.cpu.Registers[v.Name]
	case Immediate:
		return v.Value
	}

	panic("Unknown dest")
}
