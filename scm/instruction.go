package scm

import (
	"bytes"
	"fmt"
	"strings"
)

type instructionPrototype struct {
	opcode         int
	parameterCount int
	name           string
}

var prototypes map[int]instructionPrototype = make(map[int]instructionPrototype)

func registerPrototype(opcode int, parameterCount int) {
	prototypes[opcode] = instructionPrototype{opcode: opcode, parameterCount: parameterCount}
}

func applyPrototypeName(opcode int, name string) {
	prototype := prototypes[opcode]
	prototype.name = name

	prototypes[opcode] = prototype
}

// A basic SCM instruction.
// See https://gtamods.com/wiki/SCM_Instruction for compiled structure.
type Instruction struct {
	Opcode            int
	InvertReturnValue bool
	Arguments         []Value
}

func ReadInstruction(reader *bytes.Reader) Instruction {
	var opcodeBytes [2]byte
	reader.Read(opcodeBytes[:])

	opcodeMask := (uint16(opcodeBytes[1]) << 8) | uint16(opcodeBytes[0])

	instruction := Instruction{
		Opcode:            int(opcodeMask & 0x7fff),
		InvertReturnValue: opcodeMask>>0xf&1 != 0,
		Arguments:         []Value{},
	}

	prototype := prototypes[instruction.Opcode]

	for i := 0; i < prototype.parameterCount; i++ {
		instruction.Arguments = append(instruction.Arguments, ReadValue(reader))
	}

	return instruction
}

func (instruction Instruction) CodeString() string {
	parameterStrings := make([]string, len(instruction.Arguments))

	for i, argument := range instruction.Arguments {
		parameterStrings[i] = argument.CodeString()
	}

	parametersJoined := strings.Join(parameterStrings, ", ")

	nameString := fmt.Sprintf("0x%x", instruction.Opcode)

	if prototype, ok := prototypes[instruction.Opcode]; ok {
		nameString = prototype.name
	}

	result := fmt.Sprintf("%s(%s)", nameString, parametersJoined)

	if instruction.InvertReturnValue {
		result = fmt.Sprintf("!(%s)", result)
	}

	return result
}