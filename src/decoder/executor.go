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

func newX86Executor(cpu *CPU) *X86Executor {
	return &X86Executor{cpu: cpu}
}

func newX86StdoutExecutor(cpu *CPU) *X86StdoutExecutor {
	return &X86StdoutExecutor{cpu: cpu}
}

// func (e *X86Executor) Execute(instr Instruction) error {
// 	switch instr.operation {
// 	case MOV:
// 		return e.executeMOV(instr.operation)
// 	case ADD:
// 		return e.executeADD(instr)
// 	case SUB:
// 		return e.executeSUB(instr)
// 	case CMP:
// 		return e.executeCMP(instr)
// 	case JMP:
// 		return e.executeJMP(instr)
// 	default:
// 		return 	fmt.Errorf("unsupported instruction type: %v", instr.InstrType)
// 	}
// }

func (e *X86StdoutExecutor) Execute(instr Instruction) error {
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
