package main

import (
	"fmt"
	"strings"
)

const mov_inst = 0b100010
const immediate_to_reg_mov = 0b1011
const immediate_to_reg_mem = 0b1100011
const memory_to_accumulator = 0b1010000
const accumulator_to_memory = 0b1010001

const add_instr = 0
const sub_instr = 0b101
const cmp_instr = 0b111

const mov_mnemonic = "mov"
const ax = "ax"
const reg_mask byte = 0b111

var mnemonics = [4]string{mov_mnemonic, "add", "sub", "cmp"}

var fields = [16]string{
	"al", "ax", "cl", "cx", "dl", "dx", "bl", "bx",
	"ah", "sp", "ch", "bp", "dh", "si", "bh", "di",
}

var effective_addresses = [8]string{
	"bx + si", "bx + di", "bp + si", "bp + di", "si", "di", "bp", "bx",
}

func decode(bytes *[]byte) (instruction string, consumed byte) {

	switch {
	case (*bytes)[0]>>2 == mov_inst:
		instruction, consumed = decodeMov(bytes, &mnemonics[0])
	case (*bytes)[0]>>3 == add_instr:
		instruction, consumed = decodeMov(bytes, &mnemonics[1])
	case (*bytes)[0]>>3 == sub_instr:
		instruction, consumed = decodeMov(bytes, &mnemonics[2])
	case (*bytes)[0]>>3 == cmp_instr:
		instruction, consumed = decodeMov(bytes, &mnemonics[3])
	case (*bytes)[0]>>1 == immediate_to_reg_mem:
		instruction, consumed = decodeImmediateToRegMem(bytes, &mnemonics[0])
	case (*bytes)[1]>>3&reg_mask == add_instr:
		instruction, consumed = decodeImmediateToRegMem(bytes, &mnemonics[1])
	case (*bytes)[1]>>3&reg_mask == sub_instr:
		instruction, consumed = decodeImmediateToRegMem(bytes, &mnemonics[2])
	case (*bytes)[1]>>3&reg_mask == cmp_instr:
		instruction, consumed = decodeImmediateToRegMem(bytes, &mnemonics[3])
	case (*bytes)[0]>>4 == immediate_to_reg_mov:
		instruction, consumed = decodeImmediateMov(bytes)
	case (*bytes)[0]>>1 == memory_to_accumulator || (*bytes)[0]>>1 == accumulator_to_memory:
		instruction, consumed = decodeAccumulator(bytes, &mnemonics[0])
	case (*bytes)[0]>>3&reg_mask == add_instr:
		instruction, consumed = decodeAccumulator(bytes, &mnemonics[1])
	case (*bytes)[0]>>3&reg_mask == sub_instr:
		instruction, consumed = decodeAccumulator(bytes, &mnemonics[2])
	case (*bytes)[0]>>3&reg_mask == cmp_instr:
		instruction, consumed = decodeAccumulator(bytes, &mnemonics[3])
	default:
		panic("Unknown instruction")
	}

	return
}

func decodeImmediateToRegMem(bytes *[]byte, mnemonic *string) (instr string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", *mnemonic))
	w := (*bytes)[0] & 1
	consumed = 2

	rm_decoded, rm_consumed := decode_mov_rm(bytes, &w)

	dataBytes := (*bytes)[consumed+rm_consumed:]
	data, dataConsumed := decodeData(&dataBytes, &w)
	size := "byte"

	if dataConsumed == 2 {
		size = "word"
	}

	writeOperands(&builder, rm_decoded, fmt.Sprintf("%s %d", size, data))

	return builder.String(), consumed + rm_consumed + dataConsumed
}

func decodeMov(bytes *[]byte, mnemonic *string) (instruction string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", *mnemonic))
	firstByte := (*bytes)[0]

	d := (firstByte >> 1) & 1
	reg := ((*bytes)[1] >> 3) & reg_mask
	w := (*bytes)[0] & 1

	consumed = 2
	reg_decoded := decodeRegister(reg, &w)
	rm_decoded, cons := decode_mov_rm(bytes, &w)

	if d == 1 {
		writeOperands(&builder, reg_decoded, rm_decoded)
	} else {
		writeOperands(&builder, rm_decoded, reg_decoded)
	}

	return builder.String(), consumed + cons
}

func decodeImmediateMov(bytes *[]byte) (instruction string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", mov_mnemonic))

	firstByte := (*bytes)[0]
	w := (firstByte >> 3) & 1
	reg := firstByte & reg_mask

	builder.WriteString(fmt.Sprintf("%s, ", decodeRegister(reg, &w)))

	consumed = 1

	dataBytes := (*bytes)[consumed:]
	data, dataConsumed := decodeData(&dataBytes, &w)
	builder.WriteString(fmt.Sprint(data))

	return builder.String(), consumed + dataConsumed
}

func decodeAccumulator(bytes *[]byte, mnemonic *string) (instruction string, consumed byte) {
	builder := strings.Builder{}
	consumed = 1
	firstByte := (*bytes)[0]

	w := firstByte & 1

	dataBytes := (*bytes)[consumed:]
	data, dataConsumed := decodeData(&dataBytes, &w)

	dataStr := fmt.Sprint(data)

	if *mnemonic == mov_mnemonic {
		dataStr = fmt.Sprintf("[%s]", dataStr)
	}

	if (firstByte>>1)&1 == 0 {
		writeOperands(&builder, ax, dataStr)
	} else {
		writeOperands(&builder, dataStr, ax)
	}

	return builder.String(), consumed + dataConsumed
}

func decodeData(bytes *[]byte, w *byte) (data int, consumed byte) {
	data = int((*bytes)[0])

	consumed = 1

	if *w == 1 {
		data2 := int((*bytes)[1]) << 8
		data = data2 | data
		consumed++
	}

	return
}

func decode_mov_rm(bytes *[]byte, w *byte) (rm_decoded string, consumed byte) {
	mod := ((*bytes)[1] >> 6)
	rm := ((*bytes)[1]) & reg_mask

	// register
	if mod == 0b11 {
		rm_decoded = decodeRegister(rm, w)
	} else { // memory
		rm_decoded = decodeEffectiveAddress(rm)

		directAccess := mod == 0 && rm_decoded == "bp"

		var offset int

		if mod == 1 {
			offset = int((*bytes)[2])
			consumed = 1
		} else if mod == 0b10 || directAccess {
			offset = (int((*bytes)[3]) << 8) | int((*bytes)[2])
			consumed = 2
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

	return
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

func writeOperands(builder *strings.Builder, dest, src string) {
	builder.WriteString(fmt.Sprintf("%s, ", dest))
	builder.WriteString(src)
}
