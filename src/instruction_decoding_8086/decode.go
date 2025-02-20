package main

import "strings"

const mov_inst = 0b100010

var fields = [16]string{
	"al", "ax", "cl", "cx", "dl", "dx", "bl", "bx",
	"ah", "sp", "ch", "bp", "dh", "si", "bh", "di",
}

func decode(bytes *[]byte) string {
	var builder strings.Builder

	builder.WriteString(decodeInstruction(&(*bytes)[0]))

	var reg_mask byte = 0b111

	w := (*bytes)[0] & 1
	d := (*bytes)[0] >> 1 & 1

	rm := (*bytes)[1] & reg_mask
	reg := (*bytes)[1] >> 3 & reg_mask

	var src, dest = reg, rm
	if d == 1 {
		src, dest = rm, reg
	}

	builder.WriteString(decodeRegister(dest, &w))
	builder.WriteString(", ")
	builder.WriteString(decodeRegister(src, &w))

	return builder.String()
}

func decodeInstruction(firstByte *byte) string {
	opCode := *firstByte >> 2

	if opCode == mov_inst {
		return "mov "
	}

	panic("Unknown instruction")
}

func decodeRegister(register byte, w *byte) string {
	register = register<<1 | *w

	if register < byte(len(fields)) {
		return fields[register]
	}

	panic("Unknown register")
}
