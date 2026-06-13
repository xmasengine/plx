// package pir implement the PLX Intermediate Representation
package pir

import "errors"

type Operand int

const (
	OperandNone Operand = iota
	OperandInt
	OperandIdent
	OperandString
)

func (o Operand) String() string {
	switch o {
	case OperandNone:
		return "None"
	case OperandInt:
		return "Int"
	case OperandIdent:
		return "Ident"
	case OperandString:
		return "String"
	default:
		return ""
	}
}

const OperandShift = 8

type Instruction int

const (
	NOOP Instruction = 0
	NAME Instruction = 1 | Instruction(OperandIdent<<OperandShift)
	DATA Instruction = 2 | Instruction(OperandString<<OperandShift)
	IASM Instruction = 100 | Instruction(OperandString<<OperandShift)
)

func (i Instruction) String() string {
	switch i {
	case NOOP:
		return "NOOP"
	case NAME:
		return "NAME"
	case DATA:
		return "DATA"
	case IASM:
		return "IASM"
	default:
		return ""
	}
}

func (i *Instruction) MarshalText() ([]byte, error) {
	s := i.String()
	if s == "" {
		return nil, errors.New("unknown instruction")
	}
	return []byte(s), nil
}

func (i *Instruction) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "NOOP":
		*i = NOOP
	case "NAME":
		*i = NAME
	case "DATA":
		*i = DATA
	case "IASM":
		*i = IASM
	default:
		break
	}
	return errors.New("unknown instruction")
}
