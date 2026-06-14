// package pir implement the PLX Intermediate Representation
package pir

/*

PIR is an intermediate representation that models an abstract machine
with a data stack, a return stack, variables in read/write locations,
and data in read only locations.

*/
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
	// no operand instructions
	NOOP Operation = iota // no operation
	INCB                  // increment stop of stack byte
	INCW                  // increment stop of stack word
	DECB                  // increment stop of stack byte
	DECW                  // increment stop of stack word
	POPB                  // Drop byte from data stack
	POPW                  // Drop word from data stack
	DUPB                  // Duplicate byte on top of stack
	DUPW                  // Duplicate word on top of stack
	// byte operand instructions

	PSHB // Push byte literal [byte] to data stack.

	// word operand instructions
	PSHW // Push word literal [word] to data stack.
	// int operand instructions
	OUTB // Output TOS byte to port [int] (constant literal) and pop.
	OUTW // Output TOS word to port [int] (constant literal) and pop.
	OUTA // Output TOS address, length NXT bytes to port [int] (constant literal) and pop twice.
	INPB // Input byte from port [int], push to stack.
	INPW // Input word from port [int], push to stack.
	// ident operand instructions
	NAME // Name tag of next DATS, VARI, FUNC, etc instruction [ident].
	PSHT // Push address of tag [ident] to data stack.
	LABL // Define jump location [ident].
	JUMP // Jump to tag [ident] unconditionally.
	JPIF // Jump to tag [ident] it TOS is TRUE, pop stack.
	COND // Ident is one of [eq, gt, lt, etc], compare and push boolean to TOS.
	// data operand instructions
	DATS // Data String
	IASM
)

func (o Operation) String() string {
	switch o {
	case NOOP:
		return "NOOP"
	case NAME:
		return "NAME"
	case DATS:
		return "DATS"
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
	case "DATS":
		*o = DATS
	case "IASM":
		*o = IASM
	default:
		break
	}
	return errors.New("unknown Operation")
}

const (
	firstNoOp     = NOOP
	lastNoOp      = DUPW
	firstByteOp   = PSHB
	lastByteOp    = PSHB
	firstWordOp   = PSHW
	lastWordOp    = PSHW
	firstIntOp    = OUTB
	lastIntOp     = INPW
	firstIdentOp  = NAME
	lastIdentOp   = NAME
	firstStringOp = COND
	lastStringOp  = IASM
)

func (o Operation) Operand() Operand {
	switch {
	case o >= firstNoOp && o <= lastNoOp:
		return OperandNone
	case o >= firstByteOp && o <= lastByteOp:
		return OperandByte
	case o >= firstWordOp && o <= lastWordOp:
		return OperandWord
	case o >= firstIntOp && o <= lastIntOp:
		return OperandIdent
	case o >= firstIdentOp && o <= lastIdentOp:
		return OperandIdent
	case o >= firstStringOp && o <= lastStringOp:
		return OperandString
	default:
		return OperandNone
	}
}

// Instructions consist of an opcode and a single optional operand.
// This means some instructions have to look back at the previous
// instruction to get all their operands if one is not enough.
type Instruction struct {
	Operation        // Operation
	Text      string // filled in for Ident or String
	Number    int    // filled in for Byte, Word or Int
}

// Program is a list of instructions
type Program []Instruction
