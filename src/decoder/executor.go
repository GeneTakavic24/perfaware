package main

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
