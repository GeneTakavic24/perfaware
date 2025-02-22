package main

import (
	"fmt"
	"strings"
)

const mov_inst = 0b100010
const immediate_to_reg_mov = 0b1011
const mov_mnemonic = "mov"
const reg_mask byte = 0b111

var fields = [16]string{
	"al", "ax", "cl", "cx", "dl", "dx", "bl", "bx",
	"ah", "sp", "ch", "bp", "dh", "si", "bh", "di",
}

func decode(bytes *[]byte) (instruction string, consumed int) {
	var builder strings.Builder

	switch {

	case (*bytes)[0]>>2 == mov_inst:
		consumed = decodeMov(&builder, bytes, reg_mask)
	case (*bytes)[0]>>4 == immediate_to_reg_mov:
		consumed = decodeImmediateMov(&builder, bytes, reg_mask)

	default:
		panic("Unknown instruction")
	}

	return builder.String(), consumed
}

func decodeRegister(register byte, w *byte) string {
	register = register<<1 | *w

	if register < byte(len(fields)) {
		return fields[register]
	}

	panic("Unknown register")
}

func writeOperands(builder *strings.Builder, dest, src, w byte) {
	builder.WriteString(decodeRegister(dest, &w))
	builder.WriteString(", ")
	builder.WriteString(decodeRegister(src, &w))
}

func decodeMov(builder *strings.Builder, bytes *[]byte, regMask byte) int {
	builder.WriteString(mov_mnemonic + " ")

	firstByte := (*bytes)[0]
	w := firstByte & 1
	d := (firstByte >> 1) & 1
	reg := ((*bytes)[1] >> 3) & regMask
	mod := ((*bytes)[1] >> 6)
	rm := ((*bytes)[1]) & regMask

	if mod == 0b11 {
		builder.WriteString(" ")
		if d == 1 {
			writeOperands(builder, reg, rm, w)
		} else {
			writeOperands(builder, rm, reg, w)
		}
	}

	return 2
}

func decodeImmediateMov(builder *strings.Builder, bytes *[]byte, regMask byte) int {
	builder.WriteString(mov_mnemonic + " ")

	firstByte := (*bytes)[0]
	w := (firstByte >> 3) & 1
	reg := firstByte & regMask
	data1 := int((*bytes)[1])
	consumed := 2

	builder.WriteString(decodeRegister(reg, &w))
	builder.WriteString(", ")

	if w == 1 {
		data2 := int((*bytes)[2]) << 8
		data := data2 | data1
		consumed++
		builder.WriteString(fmt.Sprint(data))
	} else {
		builder.WriteString(fmt.Sprint(data1))
	}

	return consumed
}
