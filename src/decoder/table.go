package main

import (
	"fmt"
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

type EffectiveAddress struct {
	Base   Register
	Index  Register
	Offset int16
}

var effective_addresses = [8]EffectiveAddress{
	{Base: Register{"bx"}, Index: Register{"si"}},
	{Base: Register{"bx"}, Index: Register{"di"}},
	{Base: Register{"bp"}, Index: Register{"si"}},
	{Base: Register{"bp"}, Index: Register{"di"}},
	{Base: Register{"si"}},
	{Base: Register{"di"}},
	{Base: Register{"bp"}},
	{Base: Register{"bx"}},
}

func (ea EffectiveAddress) String() string {
	var b strings.Builder
	b.WriteString("[")

	hasReg := false
	if ea.Base.Name != "" {
		b.WriteString(ea.Base.Name)
		hasReg = true
	}
	if ea.Index.Name != "" {
		if hasReg {
			b.WriteString(" + ")
		}
		b.WriteString(ea.Index.Name)
		hasReg = true
	}
	if ea.Offset != 0 || !hasReg {
		if hasReg && ea.Offset > 0 {
			b.WriteString(" + ")
		} else if hasReg {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprint(ea.Offset))
	}

	b.WriteString("]")
	return b.String()
}
