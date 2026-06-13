// package pir implement the PLX Intermediate Representation
package pir

import "errors"

type Operand int

const (
	OperandNone Operand = iota
	OperandByte
	OperandWord
	OperandInt
	OperandIdent
	OperandString
)

func (o Operand) String() string {
	switch o {
	case OperandNone:
		return "None"
	case OperandByte:
		return "Byte"
	case OperandWord:
		return "Word"
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

type Operation int

const (
	NOOP Operation = 0
	NAME Operation = 1 | Operation(OperandIdent<<OperandShift)
	DATA Operation = 2 | Operation(OperandString<<OperandShift)
	IASM Operation = 100 | Operation(OperandString<<OperandShift)
)

func (o Operation) String() string {
	switch o {
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

func (o Operation) MarshalText() ([]byte, error) {
	s := o.String()
	if s == "" {
		return nil, errors.New("unknown Operation")
	}
	return []byte(s), nil
}

func (o *Operation) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "NOOP":
		*o = NOOP
	case "NAME":
		*o = NAME
	case "DATA":
		*o = DATA
	case "IASM":
		*o = IASM
	default:
		break
	}
	return errors.New("unknown Operation")
}

func (o Operation) Operand() Operand {
	return Operand(o >> OperandShift)
}

type Instruction struct {
	Operation        // Operation
	Text      string // filled in for Ident or String
	Number    int    // filled in for Byte, Word or Int
}

// Program is a list of instructions
type Program []Instruction
