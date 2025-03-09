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
	IsJump       bool
}

func newX86Executor(cpu *CPU) *X86Executor {
	return &X86Executor{cpu: cpu}
}

func (e *X86Executor) Execute(instr Instruction) error {
	prevIp := e.cpu.Registers["ip"]

	if instr.Operation == Jnz {
		offset := e.extractFrom(instr.Dest)
		jumped := e.cpu.ExecuteJump(instr.Operation, offset)

		if !jumped {
			e.cpu.Registers["ip"] = prevIp + int(instr.Consumed)
		}
	} else {
		e.cpu.Registers["ip"] = prevIp + int(instr.Consumed)
		value := e.extractFrom(instr.Src)
		if dest, ok := instr.Dest.(Register); ok {
			e.cpu.ExecuteReg(instr.Operation, dest, value)
		}
		if dest, ok := instr.Dest.(EffectiveAddress); ok {
			e.cpu.ExecuteMem(instr.Operation, dest, value)
		}
	}

	defer func() {
		if e.cpu.Registers["ip"] != prevIp {
			fmt.Printf("  ip:0x%x->0x%x", prevIp, e.cpu.Registers["ip"])
		}
	}()

	return nil
}

func (e *X86Executor) extractFrom(o Operand) int {
	switch v := o.(type) {
	case Register:
		return e.cpu.Registers[v.Name]
	case Immediate:
		return v.Value
	case EffectiveAddress:
		high, low := e.cpu.resolveEffectiveAddress(v)
		h := uint(*high) << 8
		return int(h | uint(*low))
	}

	panic("Unknown dest")
}
