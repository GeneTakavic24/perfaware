package main

import (
	"fmt"
)

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

func newX86Executor(cpu *CPU) *X86Executor {
	return &X86Executor{cpu: cpu}
}

func (e *X86Executor) Execute(instr Instruction) error {
	value := e.extractFrom(instr.Src)

	if dest, ok := instr.Dest.(Register); ok {
		e.executeToReg(instr.Operation, dest, value)
	}

	return nil
}

func (e *X86Executor) executeToReg(operation Operation, dest Register, value int) {
	if opInfo, ok := ops[operation]; ok {
		current := e.cpu.Registers[dest.Name]
		newValue := e.cpu.Execute(opInfo, dest.Name, value)
		fmt.Printf(" ; %s:%#x->%#x", dest.Name, current, newValue)
		e.cpu.printFlags()
	}
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
