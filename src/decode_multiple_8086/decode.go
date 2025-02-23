package main

import (
	"fmt"
	"strings"
)

const mov_inst = 0b100010
const immediate_to_reg_mov = 0b1011
const memory_to_accumulator = 0b1010000
const accumulator_to_memory = 0b1010001
const mov_mnemonic = "mov"
const ax = "ax"
const reg_mask byte = 0b111

var fields = [16]string{
	"al", "ax", "cl", "cx", "dl", "dx", "bl", "bx",
	"ah", "sp", "ch", "bp", "dh", "si", "bh", "di",
}

var effective_addresses = [8]string{
	"bx + si", "bx + di", "bp + si", "bp + di", "si", "di", "bp", "bx",
}

func decode(bytes *[]byte) (instruction string, consumed int) {

	switch {
	case (*bytes)[0]>>2 == mov_inst:
		instruction, consumed = decodeMov(bytes)
	case (*bytes)[0]>>4 == immediate_to_reg_mov:
		instruction, consumed = decodeImmediateMov(bytes)
	case (*bytes)[0]>>1 == memory_to_accumulator || (*bytes)[0]>>1 == accumulator_to_memory:
		instruction, consumed = decodeAccumulator(bytes)
	default:
		panic("Unknown instruction")
	}

	return
}

func decodeMov(bytes *[]byte) (instruction string, consumed int) {
	builder := strings.Builder{}
	builder.WriteString(mov_mnemonic + " ")

	firstByte := (*bytes)[0]
	w := firstByte & 1
	d := (firstByte >> 1) & 1
	reg := ((*bytes)[1] >> 3) & reg_mask
	mod := ((*bytes)[1] >> 6)
	rm := ((*bytes)[1]) & reg_mask

	var rm_decoded string

	reg_decoded := decodeRegister(reg, &w)
	consumed = 2

	// register
	if mod == 0b11 {
		rm_decoded = decodeRegister(rm, &w)
	} else { // memory
		rm_decoded = decodeEffectiveAddress(rm)

		directAccess := mod == 0 && rm_decoded == "bp"

		var offset int

		if mod == 1 {
			offset = int((*bytes)[2])
			consumed = 3
		} else if mod == 0b10 || directAccess {
			offset = (int((*bytes)[3]) << 8) | int((*bytes)[2])
			consumed = 4
		}

		if offset == 0 {
			rm_decoded = fmt.Sprintf("[%s]", rm_decoded)
		} else {
			if directAccess {
				rm_decoded = fmt.Sprintf("[%d]", offset)
			} else {
				rm_decoded = fmt.Sprintf("[%s + %d]", rm_decoded, offset)
			}
		}
	}

	if d == 1 {
		writeOperands(&builder, reg_decoded, rm_decoded)
	} else {
		writeOperands(&builder, rm_decoded, reg_decoded)
	}

	return builder.String(), consumed
}

func decodeImmediateMov(bytes *[]byte) (instruction string, consumed int) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", mov_mnemonic))

	firstByte := (*bytes)[0]
	w := (firstByte >> 3) & 1
	reg := firstByte & reg_mask

	builder.WriteString(fmt.Sprintf("%s, ", decodeRegister(reg, &w)))

	consumed = 2

	data1 := int((*bytes)[1])

	if w == 0 {
		builder.WriteString(fmt.Sprint(data1))
	} else {
		data2 := int((*bytes)[2]) << 8
		data := data2 | data1
		consumed++
		builder.WriteString(fmt.Sprint(data))
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

func decodeEffectiveAddress(rm byte) string {
	if rm < byte(len(fields)) {
		return effective_addresses[rm]
	}

	panic("Unknown register")
}

func decodeAccumulator(bytes *[]byte) (instruction string, consumed int) {
	builder := strings.Builder{}
	consumed = 2
	firstByte := (*bytes)[0]

	w := firstByte & 1
	data := int((*bytes)[1])

	if w == 1 {
		data2 := int((*bytes)[2]) << 8
		data = data2 | data
		consumed++
	}

	dataStr := fmt.Sprintf("[%d]", data)

	if (firstByte>>1)&1 == 0 {
		writeOperands(&builder, ax, dataStr)
	} else {
		writeOperands(&builder, dataStr, ax)
	}

	return builder.String(), consumed
}

func writeOperands(builder *strings.Builder, dest, src string) {
	builder.WriteString(fmt.Sprintf("%s, ", dest))
	builder.WriteString(src)
}
