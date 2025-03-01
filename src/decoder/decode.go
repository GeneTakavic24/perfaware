package main

import (
	"fmt"
	"os"
	"strings"
)

const mov_inst = 0b100010
const immediate_to_reg_mov = 0b1011
const mov_immediate_to_reg_mem = 0b1100011
const add_immediate_to_reg_mem = 0b100000
const memory_to_accumulator = 0b1010000
const accumulator_to_memory = 0b1010001

const add_instr = 0b000
const sub_reg_mem = 0b001010
const cmp_reg_mem = 0b001110
const sub_instr = 0b101
const cmp_instr = 0b111
const add_im_from_accumulator = 0b0000010
const sub_im_from_accumulator = 0b0010110
const cmp_im_from_accumulator = 0b0011110

var mov_mnemonic = "mov"

const ax = "ax"
const reg_mask byte = 0b111

type Opcode byte

const (
	JNZ    = 0b01110101
	JE     = 0b01110100
	JL     = 0b01111100
	JLE    = 0b01111110
	JB     = 0b01110010
	JBE    = 0b01110110
	JP     = 0b01111010
	JO     = 0b01110000
	JS     = 0b01111000
	JNL    = 0b01111101
	JG     = 0b01111111
	JNB    = 0b01110011
	JA     = 0b01110111
	JNP    = 0b01111011
	JNO    = 0b01110001
	JNS    = 0b01111001
	LOOP   = 0b11100010
	LOOPZ  = 0b11100001
	LOOPNZ = 0b11100000
	JCXZ   = 0b11100011
)

var opcodeNames = map[Opcode]string{
	JNZ:    "jnz",
	JE:     "je",
	JL:     "jl",
	JLE:    "jle",
	JB:     "jb",
	JBE:    "jbe",
	JP:     "jp",
	JO:     "jo",
	JS:     "js",
	JNL:    "jnl",
	JG:     "jg",
	JNB:    "jnb",
	JA:     "ja",
	JNP:    "jnp",
	JNO:    "jno",
	JNS:    "jns",
	LOOP:   "loop",
	LOOPZ:  "loopz",
	LOOPNZ: "loopnz",
	JCXZ:   "jcxz",
}

var fields = [16]string{
	"al", "ax", "cl", "cx", "dl", "dx", "bl", "bx",
	"ah", "sp", "ch", "bp", "dh", "si", "bh", "di",
}

var effective_addresses = [8]string{
	"bx + si", "bx + di", "bp + si", "bp + di", "si", "di", "bp", "bx",
}

func getMnemonic(mnemonic byte) string {
	switch mnemonic {
	case add_instr:
		return "add"
	case sub_instr:
		return "sub"
	case cmp_instr:
		return "cmp"
	}
	panic("Unknown mnemonic")
}

func decode(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("; %s disassembly:\n", filePath)
	fmt.Println("bits 16")
	fmt.Println()

	for i := 0; i < len(bytes); {
		end := min(i+6, len(bytes))
		instr, consumed := parseInstruction(bytes[i:end])
		fmt.Println(instr)
		i += int(consumed)
	}
}

func parseInstruction(bytes []byte) (instruction string, consumed byte) {
	if _, ok := opcodeNames[Opcode(bytes[0])]; ok {
		return decodeJmp(bytes)
	}

	switch {
	case bytes[0]>>4 == immediate_to_reg_mov:
		return decodeImmediateMov(bytes)
	}

	switch bytes[0] >> 2 {
	case mov_inst:
		return decodeRegMem(bytes, &mov_mnemonic)
	case add_instr, sub_reg_mem, cmp_reg_mem:
		return decodeRegMemWrapper(bytes)
	case add_immediate_to_reg_mem, sub_instr, cmp_instr:
		return decodeImmediateToRegMemWrapper(bytes)
	}

	switch bytes[0] >> 1 {
	case mov_immediate_to_reg_mem:
		return decodeImmediateToRegMem(bytes, &mov_mnemonic)
	case memory_to_accumulator, accumulator_to_memory:
		return decodeAccumulator(bytes, &mov_mnemonic)
	case add_im_from_accumulator, sub_im_from_accumulator, cmp_im_from_accumulator:
		return decodeAccumulatorWrapper(bytes)
	}

	panic("Unknown instruction")
}

func decodeJmp(bytes []byte) (instr string, consumed byte) {
	builder := strings.Builder{}
	mnemonic := opcodeNames[Opcode(bytes[0])]
	builder.WriteString(fmt.Sprintf("%s %d", mnemonic, int8(bytes[1])))
	return builder.String(), 2
}

func decodeImmediateToRegMemWrapper(bytes []byte) (instr string, consumed byte) {
	op := (bytes[1] >> 3) & reg_mask
	mnemonic := getMnemonic(op)
	return decodeImmediateToRegMem(bytes, &mnemonic)
}

func decodeImmediateToRegMem(bytes []byte, mnemonic *string) (instr string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", *mnemonic))
	w := bytes[0] & 1
	consumed = 2
	mod := bytes[1] >> 6

	rm_decoded, rm_consumed := decode_mov_rm(bytes, &w, &mod)

	dataBytes := bytes[consumed+rm_consumed:]
	var data uint16
	var dataConsumed byte

	if *mnemonic == mov_mnemonic {
		data, dataConsumed = decodeData(dataBytes, &w)
	} else {
		s := (bytes[0] >> 1) & 1
		w = (^s & 1) & w
		data, dataConsumed = decodeData(dataBytes, &w)
	}

	dataStr := fmt.Sprint(data)

	if mod != 0b11 {
		size := "byte"

		if dataConsumed == 2 {
			size = "word"
		}

		dataStr = fmt.Sprintf("%s %s", size, dataStr)
	}

	writeOperands(&builder, rm_decoded, dataStr)

	return builder.String(), consumed + rm_consumed + dataConsumed
}

func decodeRegMemWrapper(bytes []byte) (instruction string, consumed byte) {
	op := (bytes[0] >> 3) & reg_mask
	mnemonic := getMnemonic(op)
	return decodeRegMem(bytes, &mnemonic)
}

func decodeRegMem(bytes []byte, mnemonic *string) (instruction string, consumed byte) {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s ", *mnemonic))
	firstByte := bytes[0]

	d := (firstByte >> 1) & 1
	reg := (bytes[1] >> 3) & reg_mask
	w := bytes[0] & 1
	mod := bytes[1] >> 6

	consumed = 2
	reg_decoded := decodeRegister(reg, &w)
	rm_decoded, cons := decode_mov_rm(bytes, &w, &mod)

	if d == 1 {
		writeOperands(&builder, reg_decoded, rm_decoded)
	} else {
		writeOperands(&builder, rm_decoded, reg_decoded)
	}

	return builder.String(), consumed + cons
}

func decodeImmediateMov(bytes []byte) (instruction string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", mov_mnemonic))

	firstByte := bytes[0]
	w := (firstByte >> 3) & 1
	reg := firstByte & reg_mask

	builder.WriteString(fmt.Sprintf("%s, ", decodeRegister(reg, &w)))

	consumed = 1

	dataBytes := bytes[consumed:]
	data, dataConsumed := decodeData(dataBytes, &w)
	builder.WriteString(fmt.Sprint(data))

	return builder.String(), consumed + dataConsumed
}

func decodeAccumulatorWrapper(bytes []byte) (instruction string, consumed byte) {
	op := bytes[0] >> 3 & reg_mask
	mnemonic := getMnemonic(op)
	return decodeAccumulator(bytes, &mnemonic)
}

func decodeAccumulator(bytes []byte, mnemonic *string) (instruction string, consumed byte) {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s ", *mnemonic))

	consumed = 1
	firstByte := bytes[0]

	w := firstByte & 1

	dataBytes := bytes[consumed:]
	data, dataConsumed := decodeData(dataBytes, &w)

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

func decodeData(bytes []byte, w *byte) (data uint16, consumed byte) {
	data = uint16(bytes[0])

	consumed = 1

	if *w == 1 {
		data2 := uint16(bytes[1]) << 8
		data = data2 | data
		consumed++
	}

	return
}

func decode_mov_rm(bytes []byte, w, mod *byte) (rm_decoded string, consumed byte) {
	rm := bytes[1] & reg_mask

	// register
	if *mod == 0b11 {
		rm_decoded = decodeRegister(rm, w)
	} else { // memory
		rm_decoded = decodeEffectiveAddress(rm)

		directAccess := *mod == 0 && rm_decoded == "bp"

		var offset int16

		if *mod == 1 {
			offset = int16(int8(bytes[2]))
			consumed = 1
		} else if *mod == 0b10 || directAccess {
			offset = int16(bytes[2]) | (int16(bytes[3]) << 8)
			consumed = 2
		}

		if offset == 0 {
			rm_decoded = fmt.Sprintf("[%s]", rm_decoded)
		} else {
			if directAccess {
				rm_decoded = fmt.Sprintf("[%d]", offset)
			} else {
				sign := "+"
				if offset < 0 {
					sign = "-"
					offset = -offset
				}
				rm_decoded = fmt.Sprintf("[%s %s %d]", rm_decoded, sign, offset)
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
