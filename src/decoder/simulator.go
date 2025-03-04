package main

import (
	"fmt"
	"os"
)

func simulate(filePath string) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	cpu := NewCPU(10)

	executor := newX86Executor(cpu)

	fmt.Printf("; %s disassembly:\n", filePath)
	fmt.Println("bits 16")
	fmt.Println()

	for i := 0; i < len(bytes); {
		end := min(i+6, len(bytes))
		instr := parseInstruction(bytes[i:end])
		Stdout(instr)
		executor.Execute(instr)
		fmt.Println()
		i += int(instr.Consumed)
	}

	cpu.PrintCPU()
}

func parseInstruction(bytes []byte) Instruction {
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

func decodeJmp(bytes []byte) Instruction {
	mnemonic := opcodeNames[Opcode(bytes[0])]

	return Instruction{
		Operation: Operation(mnemonic),
		Dest:      Immediate{Value: int(int8(bytes[1]))},
		Consumed:  2,
	}
}

func decodeImmediateToRegMemWrapper(bytes []byte) Instruction {
	op := (bytes[1] >> 3) & reg_mask
	mnemonic := getMnemonic(op)
	return decodeImmediateToRegMem(bytes, &mnemonic)
}

func decodeImmediateToRegMem(bytes []byte, mnemonic *string) Instruction {
	w := bytes[0] & 1
	mod := bytes[1] >> 6

	instruction := Instruction{
		Operation: Operation(*mnemonic),
		Consumed:  2,
	}

	operand, rm_consumed := decode_mov_rm(bytes, &w, &mod)
	instruction.Consumed += rm_consumed

	dataBytes := bytes[instruction.Consumed:]
	var data uint16
	var dataConsumed byte

	if *mnemonic == mov_mnemonic {
		data, dataConsumed = decodeData(dataBytes, &w)
	} else {
		s := (bytes[0] >> 1) & 1
		w = (^s & 1) & w
		data, dataConsumed = decodeData(dataBytes, &w)
	}

	immediateData := Immediate{Value: int(data)}

	if mod != 0b11 {
		immediateData.ExplicitSize = "byte"

		if dataConsumed == 2 {
			immediateData.ExplicitSize = "word"
		}
	}

	instruction.Dest = operand
	instruction.Src = immediateData
	instruction.Consumed += dataConsumed

	return instruction
}

func decodeRegMemWrapper(bytes []byte) Instruction {
	op := (bytes[0] >> 3) & reg_mask
	mnemonic := getMnemonic(op)
	return decodeRegMem(bytes, &mnemonic)
}

func decodeRegMem(bytes []byte, mnemonic *string) Instruction {
	d := (bytes[0] >> 1) & 1
	reg := (bytes[1] >> 3) & reg_mask
	w := bytes[0] & 1
	mod := bytes[1] >> 6

	instruction := Instruction{
		Operation: Operation(*mnemonic),
		Consumed:  2,
	}

	reg_decoded := decodeRegister(reg, &w)
	rm_decoded, cons := decode_mov_rm(bytes, &w, &mod)
	instruction.Consumed += cons

	if d == 1 {
		instruction.Dest = reg_decoded
		instruction.Src = rm_decoded
	} else {
		instruction.Dest = rm_decoded
		instruction.Src = reg_decoded
	}

	return instruction
}

func decodeImmediateMov(bytes []byte) Instruction {
	firstByte := bytes[0]
	w := (firstByte >> 3) & 1
	reg := firstByte & reg_mask

	instruction := Instruction{
		Operation: Operation(mov_mnemonic),
		Consumed:  1,
	}

	dataBytes := bytes[instruction.Consumed:]
	data, dataConsumed := decodeData(dataBytes, &w)
	instruction.Consumed += dataConsumed
	instruction.Dest = decodeRegister(reg, &w)
	instruction.Src = Immediate{Value: int(data)}

	return instruction
}

func decodeAccumulatorWrapper(bytes []byte) Instruction {
	op := bytes[0] >> 3 & reg_mask
	mnemonic := getMnemonic(op)
	return decodeAccumulator(bytes, &mnemonic)
}

func decodeAccumulator(bytes []byte, mnemonic *string) Instruction {
	instruction := Instruction{
		Operation: Operation(*mnemonic),
		Consumed:  1,
	}
	firstByte := bytes[0]

	w := firstByte & 1

	dataBytes := bytes[instruction.Consumed:]
	data, dataConsumed := decodeData(dataBytes, &w)
	instruction.Consumed += dataConsumed

	register := Register{Name: ax}
	var operand Operand = Immediate{Value: int(data)}

	if *mnemonic == mov_mnemonic {
		operand = EffectiveAddress{
			Offset: int16(data),
		}
	}

	if (firstByte>>1)&1 == 0 {
		instruction.Dest = register
		instruction.Src = operand
	} else {
		instruction.Dest = operand
		instruction.Src = register
	}

	return instruction
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

func decode_mov_rm(bytes []byte, w, mod *byte) (operand Operand, consumed byte) {
	rm := bytes[1] & reg_mask

	// register
	if *mod == 0b11 {
		operand = decodeRegister(rm, w)
	} else { // memory
		effectiveAddress := decodeEffectiveAddress(rm)

		directAccess := *mod == 0 && rm == 0b110

		if *mod == 1 {
			effectiveAddress.Offset = int16(int8(bytes[2]))
			consumed = 1
		} else if *mod == 0b10 || directAccess {
			effectiveAddress.Offset = int16(bytes[2]) | (int16(bytes[3]) << 8)
			consumed = 2
		}

		if effectiveAddress.Offset != 0 && directAccess {
			effectiveAddress.Base = Register{""}
			effectiveAddress.Index = Register{""}
		}

		operand = effectiveAddress
	}

	return
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

func decodeRegister(register byte, w *byte) Register {
	register = register<<1 | *w

	if register < byte(len(fields)) {
		return Register{fields[register]}
	}

	panic("Unknown register")
}

func decodeEffectiveAddress(rm byte) EffectiveAddress {
	if rm < byte(len(fields)) {
		return effective_addresses[rm]
	}

	panic("Unknown register")
}
