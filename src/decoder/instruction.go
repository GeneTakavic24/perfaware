package main

import (
	"fmt"
	"strings"
)

type Operand interface {
	String() string
}

type Instruction struct {
	Operation Operation
	Src       Operand
	Dest      Operand
	Consumed  byte
}

type Operation string

const (
	MOV Operation = "mov"
	ADD Operation = "add"
	SUB Operation = "sub"
)

type Register struct {
	Name string
}

type Immediate struct {
	Value        int
	ExplicitSize string
}

func (r Register) String() string {
	return r.Name
}

func (r Immediate) String() string {
	var b strings.Builder

	if r.ExplicitSize != "" {
		b.WriteString(r.ExplicitSize)
		b.WriteRune(' ')
	}

	b.WriteString(fmt.Sprintf("%d", r.Value))

	return b.String()
}
