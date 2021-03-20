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

	err = LoadDumped("data/prototypes.scmpt")

	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(codeBytes)

	instructions := make([]Instruction, 0)

	for reader.Len() != 0 {
		instruction := ReadInstruction(reader)
		instructions = append(instructions, instruction)

		if instruction.Opcode == 0 {
			continue
		}

		println(instruction.CodeString())
	}
}
