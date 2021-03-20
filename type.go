package scm

import (
	"fmt"
)

type ConcreteType byte

const (
	ConcreteEndOfArguments        ConcreteType = 0x00
	ConcreteSigned32              ConcreteType = 0x01
	ConcreteGlobal32              ConcreteType = 0x02
	ConcreteLocal32               ConcreteType = 0x03
	ConcreteSigned8               ConcreteType = 0x04
	ConcreteSigned16              ConcreteType = 0x05
	ConcreteFloat32               ConcreteType = 0x06
	ConcreteGlobal32Element       ConcreteType = 0x07
	ConcreteLocal32Element        ConcreteType = 0x08
	ConcreteString8               ConcreteType = 0x09
	ConcreteGlobalString8         ConcreteType = 0x0a
	ConcreteLocalString8          ConcreteType = 0x0b
	ConcreteGlobalString8Element  ConcreteType = 0x0c
	ConcreteLocalString8Element   ConcreteType = 0x0d
	ConcreteVariableString        ConcreteType = 0x0e
	ConcreteString16              ConcreteType = 0x0f
	ConcreteGlobalString16        ConcreteType = 0x10
	ConcreteLocalString16         ConcreteType = 0x11
	ConcreteGlobalString16Element ConcreteType = 0x12
	ConcreteLocalString16Element  ConcreteType = 0x13
)

func (concreteType ConcreteType) ValueLength() int {
	switch concreteType {
	case ConcreteEndOfArguments:
		return 0

	case ConcreteSigned32, ConcreteFloat32:
		return 4

	case ConcreteGlobal32,
		ConcreteLocalString16,
		ConcreteGlobalString16,
		ConcreteLocalString8,
		ConcreteGlobalString8,
		ConcreteSigned16,
		ConcreteLocal32:
		return 2

	case ConcreteSigned8:
		return 1

	case ConcreteGlobal32Element,
		ConcreteLocalString16Element,
		ConcreteGlobalString16Element,
		ConcreteLocalString8Element,
		ConcreteGlobalString8Element,
		ConcreteLocal32Element:
		return 6

	case ConcreteString8:
		return 8

	case ConcreteVariableString:
		return -1

	case ConcreteString16:
		return 16
	}

	return 0
}

type AbstractType byte

const (
	AbstractNil                 AbstractType = iota
	AbstractString              AbstractType = iota
	AbstractEnd                 AbstractType = iota
	AbstractGlobal32            AbstractType = iota
	AbstractLocal32             AbstractType = iota
	AbstractInteger             AbstractType = iota
	AbstractFloat               AbstractType = iota
	AbstractGlobal32Element     AbstractType = iota
	AbstractLocal32Element      AbstractType = iota
	AbstractGlobalString        AbstractType = iota
	AbstractLocalString         AbstractType = iota
	AbstractGlobalStringElement AbstractType = iota
	AbstractLocalStringElement  AbstractType = iota
)

type DataType struct {
	Concrete ConcreteType
	Abstract AbstractType
}

func (dataType DataType) IsAbstract(abstract AbstractType) bool {
	return dataType.Abstract == abstract
}

func (dataType DataType) IsLocal() bool {
	switch dataType.Abstract {
	case AbstractLocal32, AbstractLocalString:
		return true
	}

	return false
}

func (dataType DataType) IsGlobal() bool {
	switch dataType.Abstract {
	case AbstractGlobal32, AbstractGlobalString:
		return true
	}

	return false
}

func (dataType DataType) IsVariable() bool {
	return dataType.IsGlobal() || dataType.IsLocal()
}

func (dataType DataType) IsArrayElement() bool {
	switch dataType.Abstract {
	case AbstractGlobal32Element:
		return true
	case AbstractLocal32Element:
		return true
	case AbstractGlobalStringElement:
		return true
	case AbstractLocalStringElement:
		return true
	}

	return false
}

func (dataType DataType) IsConcrete(concrete ConcreteType) bool {
	return dataType.Concrete == concrete
}

func (concreteType ConcreteType) Lift() DataType {
	var abstractType AbstractType

	switch concreteType {
	case ConcreteString8, ConcreteString16, ConcreteVariableString:
		abstractType = AbstractString

	case ConcreteEndOfArguments:
		abstractType = AbstractEnd

	case ConcreteGlobal32:
		abstractType = AbstractGlobal32

	case ConcreteLocal32:
		abstractType = AbstractLocal32

	case ConcreteSigned8, ConcreteSigned16, ConcreteSigned32:
		abstractType = AbstractInteger

	case ConcreteFloat32:
		abstractType = AbstractFloat

	case ConcreteGlobal32Element:
		abstractType = AbstractGlobal32Element

	case ConcreteLocal32Element:
		abstractType = AbstractLocal32Element

	case ConcreteGlobalString8, ConcreteGlobalString16:
		abstractType = AbstractGlobalString

	case ConcreteLocalString8, ConcreteLocalString16:
		abstractType = AbstractLocalString

	case ConcreteGlobalString8Element, ConcreteGlobalString16Element:
		abstractType = AbstractGlobalStringElement

	case ConcreteLocalString8Element, ConcreteLocalString16Element:
		abstractType = AbstractLocalStringElement

	default:
		panic(fmt.Errorf("invalid concrete type: %x", concreteType))
	}

	return DataType{Concrete: concreteType, Abstract: abstractType}
}
