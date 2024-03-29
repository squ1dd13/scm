package scm

import (
	_ "embed"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readLines(path string) ([]string, error) {
	fileBytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(fileBytes), "\n"), nil
}

// Load parameter counts from an SASCM.ini file.
func LoadPrototypes(path string) error {
	lines, err := readLines(path)

	if err != nil {
		return err
	}

	for _, dirtyLine := range lines {
		line := strings.TrimSpace(dirtyLine)

		if len(line) == 0 || line[0] == ';' {
			continue
		}

		equalsIndex := strings.IndexRune(line, '=')

		if equalsIndex < 0 {
			continue
		}

		opcode, err := strconv.ParseInt(line[:equalsIndex], 16, 16)

		if err != nil {
			println("Unable to parse opcode from line: '" + line + "'.")
			continue
		}

		commaIndex := strings.IndexRune(line, ',')

		if commaIndex < 0 {
			continue
		}

		count, err := strconv.Atoi(line[equalsIndex+1 : commaIndex])

		if err != nil {
			println("Unable to parse count from line: '" + line + "'.")
			continue
		}

		// We only read the opcodes and parameter counts from this file.
		registerPrototype(int(opcode), count)
	}

	return nil
}

// Load prototype names from a commands.def file.
func LoadPrototypeNames(path string) error {
	lines, err := readLines(path)

	if err != nil {
		return err
	}

	for _, dirtyLine := range lines {
		commentIndex := strings.IndexRune(dirtyLine, ';')

		if -1 < commentIndex {
			dirtyLine = dirtyLine[:commentIndex]
		}

		line := strings.TrimSpace(dirtyLine)

		if len(line) == 0 {
			continue
		}

		opcodeHex := line[:4]
		name := strings.TrimSpace(line[4:])

		bracketIndex := strings.IndexRune(name, '(')

		if -1 < bracketIndex {
			name = name[:bracketIndex]
		}

		opcode, err := strconv.ParseInt(opcodeHex, 16, 16)

		if err != nil {
			println("Unable to parse opcode from line (name is '" + name + "'): '" + line + "'.")
			continue
		}

		applyPrototypeName(int(opcode), name)
	}

	return nil
}

//go:embed data/prototypes.scmpt
var prototypeFileString string

// Load all prototypes from a prototypes.scmpt file. Use this once the
// SASCM.ini and commands.def files have been loaded and the prototypes
// dumped.
func LoadDumped() error {
	lines := strings.Split(prototypeFileString, "\n")

	for _, dirtyLine := range lines {
		commentIndex := strings.IndexRune(dirtyLine, ';')

		if -1 < commentIndex {
			dirtyLine = dirtyLine[:commentIndex]
		}

		line := strings.TrimSpace(dirtyLine)

		if len(line) == 0 {
			continue
		}

		var opcode int
		var invokeType string
		var name string
		var parameterCount int

		read, err := fmt.Sscanf(line, "%x (%s %d) %s", &opcode, &invokeType, &parameterCount, &name)

		if err != nil || read != 4 {
			continue
		}

		addPrototype(opcode, name, parameterCount, invokeType == "oper")
	}

	return nil
}

// Dump the prototype map. Use this after loading SASCM.ini and commands.def to
// produce a .scmpt file at `path`.
func DumpPrototypes(path string) error {
	builder := strings.Builder{}

	prototypeSlice := make([]instructionPrototype, 0, len(prototypes))

	for opcode, prototype := range prototypes {
		prototype.opcode = opcode
		prototypeSlice = append(prototypeSlice, prototype)
	}

	sort.Slice(prototypeSlice, func(i, j int) bool {
		return prototypeSlice[i].opcode < prototypeSlice[j].opcode
	})

	for _, prototype := range prototypeSlice {
		invokeType := "func"

		if prototype.isOperator {
			invokeType = "oper"
		}

		fmt.Fprintf(&builder, "%04x (%s %d) %s\n", prototype.opcode, invokeType, prototype.parameterCount, prototype.name)
	}

	return os.WriteFile(path, []byte(builder.String()), 0755)
}

func init() {
	err := LoadDumped()

	if err != nil {
		panic(err)
	}
}
