package main

import (
	"fmt"
	"strings"
)

type Executor interface {
	Execute(instr Instruction) error
}

type X86Executor struct {
	cpu *CPU
}

type X86StdoutExecutor struct {
	cpu *CPU
}

var ops = map[Operation]func(current, value int) int{
	MOV: func(_, v int) int { return v },
	ADD: func(c, v int) int { return c + v },
	SUB: func(c, v int) int { return c - v },
}

func newX86Executor(cpu *CPU) *X86Executor {
	return &X86Executor{cpu: cpu}
}

func newX86StdoutExecutor(cpu *CPU) *X86StdoutExecutor {
	return &X86StdoutExecutor{cpu: cpu}
}

func (e *X86Executor) Execute(instr Instruction) error {
	value := e.extractFrom(instr.Src)

	if dest, ok := instr.Dest.(Register); ok {
		e.executeToReg(instr.Operation, dest, value)
	}

	return nil
}

func (e *X86Executor) executeToReg(operation Operation, dest Register, value int) {
	if opFunc, ok := ops[operation]; ok {
		current := e.cpu.Registers[dest.Name]
		newValue := opFunc(current, value)
		e.cpu.Registers[dest.Name] = newValue
		fmt.Printf(" ; %s:%#x->%#x", dest.Name, current, newValue)
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

func Stdout(instr Instruction) error {
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
	return nil
}
