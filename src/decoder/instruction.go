package main

type Operand any

type Instruction struct {
	Operation Operation
	Src       Operand
	Dest      Operand
	Consumed  byte
}

type Operation string

const (
	MOV = iota
	ADD
	SUB
	CMP
	JMP
)

type Register struct {
	Name string
}

type Immediate struct {
	Value        int
	ExplicitSize string
}
