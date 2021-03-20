package scm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type ElementType byte

const (
	ElementTypeInt      ElementType = iota
	ElementTypeFloat    ElementType = iota
	ElementTypeString8  ElementType = iota
	ElementTypeString16 ElementType = iota
)

type ArrayAccess struct {
	Element             ElementType
	FirstVariableOffset uint16
	IndexVariableOffset uint16
	InfoMask            byte
}

type Value struct {
	Type DataType

	integer *int64
	float   *float32
	str     *string
	array   *ArrayAccess
}

func intFromBytes(bites []byte) int64 {
	reader := bytes.NewReader(bites)

	switch len(bites) {
	case 1:
		{
			var value int8
			binary.Read(reader, binary.LittleEndian, &value)

			return int64(value)
		}

	case 2:
		{
			var value int16
			binary.Read(reader, binary.LittleEndian, &value)

			return int64(value)
		}

	case 4:
		{
			var value int32
			binary.Read(reader, binary.LittleEndian, &value)

			return int64(value)
		}
	}

	panic(errors.New("invalid count"))
}

func ReadValue(reader *bytes.Reader) Value {
	typeByte, _ := reader.ReadByte()
	dataType := ConcreteType(typeByte).Lift()

	length := dataType.Concrete.ValueLength()

	if dataType.IsConcrete(ConcreteVariableString) {
		lengthByte, _ := reader.ReadByte()
		length = int(lengthByte)
	}

	buffer := make([]byte, length)
	reader.Read(buffer)
	bufferReader := bytes.NewBuffer(buffer)

	if dataType.IsAbstract(AbstractInteger) || dataType.IsVariable() {
		fromBytes := intFromBytes(buffer)
		return Value{Type: dataType, integer: &fromBytes}
	}

	if dataType.IsAbstract(AbstractFloat) {
		var floating float32
		err := binary.Read(bufferReader, binary.LittleEndian, &floating)

		if err != nil {
			panic(err)
		}

		return Value{Type: dataType, float: &floating}
	}

	if dataType.IsAbstract(AbstractString) {
		str := string(buffer)

		if nullIndex := strings.IndexByte(str, 0); -1 < nullIndex {
			str = str[0:nullIndex]
		}

		return Value{Type: dataType, str: &str}
	}

	if dataType.IsArrayElement() {
		var arrayAccess ArrayAccess

		err := binary.Read(bufferReader, binary.LittleEndian, &arrayAccess)

		if err != nil {
			panic(err)
		}

		return Value{Type: dataType, array: &arrayAccess}
	}

	return Value{}
}

func (value Value) CodeString() string {
	if value.array != nil {
		return fmt.Sprintf("0x%x[*0x%x]", value.array.FirstVariableOffset, value.array.IndexVariableOffset)
	}

	if value.float != nil {
		return fmt.Sprint(*value.float)
	}

	if value.integer != nil {
		prefix := ""

		if value.Type.IsLocal() {
			prefix = "local_"
		} else if value.Type.IsGlobal() {
			prefix = "global_"
		}

		return prefix + fmt.Sprint(*value.integer)
	}

	if value.str != nil {
		return fmt.Sprintf("\"%s\"", *value.str)
	}

	panic(errors.New("unable to produce code string"))
}
