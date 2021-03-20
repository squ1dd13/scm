package scm

import (
	"bytes"
	"os"
	"testing"
)

func TestSCM(t *testing.T) {
	codeBytes, err := os.ReadFile(os.Args[1])

	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(codeBytes)

	instructions := make([]Instruction, 0)

	for reader.Len() != 0 {
		instruction := ReadInstruction(reader)

		if instruction == nil {
			println("Bad instruction, stopping.")
			break
		}

		instructions = append(instructions, *instruction)

		if instruction.Opcode == 0 {
			continue
		}

		println(instruction.CodeString())
	}
}
