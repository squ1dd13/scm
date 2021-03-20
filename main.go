package main

import (
	"bytes"
	"gta_scm/scm"
	"os"
)

func main() {
	codeBytes, err := os.ReadFile("/media/squ1dd13/BarraCuda/macOS/CLionProjects/gtasm/GTA Scripts/shopper.scm")

	if err != nil {
		panic(err)
	}

	// loadError := scm.LoadPrototypes("/home/squ1dd13/Documents/Projects/Java/MSD/SASCM.ini")

	// if loadError != nil {
	// 	panic(loadError)
	// }

	// loadError = scm.LoadPrototypeNames("/home/squ1dd13/Documents/Projects/Java/MSD/commands.ini")

	// if loadError != nil {
	// 	panic(loadError)
	// }

	// dumpError := scm.DumpPrototypes("/home/squ1dd13/go/src/scm/data/prototypes.scmpt")

	// if dumpError != nil {
	// 	panic(dumpError)
	// }

	scm.LoadDumped("/home/squ1dd13/go/src/scm/data/prototypes.scmpt")

	reader := bytes.NewReader(codeBytes)

	instructions := make([]scm.Instruction, 0)

	for reader.Len() != 0 {
		instruction := scm.ReadInstruction(reader)
		instructions = append(instructions, instruction)

		if instruction.Opcode == 0 {
			continue
		}

		println(instruction.CodeString())
	}
}
