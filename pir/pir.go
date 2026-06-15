// package pir implement the PLX Intermediate Representation
package pir

/*

PIR is an intermediate representation that models an abstract machine
with a data stack, a return stack, variables in read/write locations,
and data in read only locations.

TOS = top of stack, NXT is the next of stack, SP3 is the third of stack, etc.

For operations with 1 result and 0 inputs: push ; OP -> TOS
For operations with 1 result and 1 inputs: TOS -> OP -> TOS
For operations with 1 result and 2 inputs: NXT, TOS -> OP -> TOS
For operations with 1 result and 3 inputs: SP3, NXT, TOS -> OP -> TOS

The stack contains only 16 bits values. 8 bits values have the high byte 0.
This makes stack dicipline easier.

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
	PSHA // Push address of [ident] to data stack.
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

	ADDB
	ADDW
	SUBB
	SUBW
	ANDB
	ANDW
	BORB
	BORW
	XORB
	XORW
	SHLB
	SHLW
	SHRB
	SHRW
	// could also support rolls and arithmetic shifts
)

// This is only used by the parser to skip empty lines.
const SKIP Operation = -1

// This is only used by the parser to indicate the end of file.
const DONE Operation = -2

func (o Operation) String() string {
	switch o {
	case NOOP:
		return "NOOP"
	case INCB:
		return "INCB"
	case INCW:
		return "INCW"
	case DECB:
		return "DECB"
	case DECW:
		return "DECW"
	case POPB:
		return "POPB"
	case POPW:
		return "POPW"
	case DUPB:
		return "DUPB"
	case DUPW:
		return "DUPW"
	case PSHB:
		return "PSHB"
	case PSHW:
		return "PSHW"
	case PSHA:
		return "PSHA"
	case OUTB:
		return "OUTB"
	case OUTW:
		return "OUTW"
	case OUTA:
		return "OUTA"
	case INPB:
		return "INPB"
	case INPW:
		return "INPW"
	case NAME:
		return "NAME"
	case PSHT:
		return "PSHT"
	case LABL:
		return "LABL"
	case JUMP:
		return "JUMP"
	case JPIF:
		return "JPIF"
	case COND:
		return "COND"
	case DATS:
		return "DATS"
	case IASM:
		return "IASM"
	case ADDB:
		return "ADDB"
	case ADDW:
		return "ADDW"
	case SUBB:
		return "SUBB"
	case SUBW:
		return "SUBW"
	case ANDB:
		return "ANDB"
	case ANDW:
		return "ANDW"
	case BORB:
		return "BORB"
	case BORW:
		return "BORW"
	case XORB:
		return "XORB"
	case XORW:
		return "XORW"
	case SHLB:
		return "SHLB"
	case SHLW:
		return "SHLW"
	case SHRB:
		return "SHRB"
	case SHRW:
		return "SHRW"
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
	case "INCB":
		*o = INCB
	case "INCW":
		*o = INCW
	case "DECB":
		*o = DECB
	case "DECW":
		*o = DECW
	case "POPB":
		*o = POPB
	case "POPW":
		*o = POPW
	case "DUPB":
		*o = DUPB
	case "DUPW":
		*o = DUPW
	case "PSHB":
		*o = PSHB
	case "PSHW":
		*o = PSHW
	case "PSHA":
		*o = PSHA
	case "OUTB":
		*o = OUTB
	case "OUTW":
		*o = OUTW
	case "OUTA":
		*o = OUTA
	case "INPB":
		*o = INPB
	case "INPW":
		*o = INPW
	case "NAME":
		*o = NAME
	case "PSHT":
		*o = PSHT
	case "LABL":
		*o = LABL
	case "JUMP":
		*o = JUMP
	case "JPIF":
		*o = JPIF
	case "COND":
		*o = COND
	case "DATS":
		*o = DATS
	case "IASM":
		*o = IASM
	case "ADDB":
		*o = ADDB
	case "ADDW":
		*o = ADDW
	case "SUBB":
		*o = SUBB
	case "SUBW":
		*o = SUBW
	case "ANDB":
		*o = ANDB
	case "ANDW":
		*o = ANDW
	case "BORB":
		*o = BORB
	case "BORW":
		*o = BORW
	case "XORB":
		*o = XORB
	case "XORW":
		*o = XORW
	case "SHLB":
		*o = SHLB
	case "SHLW":
		*o = SHLW
	case "SHRB":
		*o = SHRB
	case "SHRW":
		*o = SHRW
	default:
		return errors.New("unknown operation: " + s)
	}
	return nil
}

func (o Operation) Operand() Operand {
	switch o {
	case NOOP:
		return OperandNone
	case INCB:
		return OperandNone
	case INCW:
		return OperandNone
	case DECB:
		return OperandNone
	case DECW:
		return OperandNone
	case POPB:
		return OperandNone
	case POPW:
		return OperandNone
	case DUPB:
		return OperandNone
	case DUPW:
		return OperandNone
	case PSHB:
		return OperandByte
	case PSHW:
		return OperandWord
	case PSHA:
		return OperandIdent
	case OUTB:
		return OperandInt
	case OUTW:
		return OperandInt
	case OUTA:
		return OperandInt
	case INPB:
		return OperandInt
	case INPW:
		return OperandInt
	case NAME:
		return OperandIdent
	case PSHT:
		return OperandIdent
	case LABL:
		return OperandIdent
	case JUMP:
		return OperandIdent
	case JPIF:
		return OperandIdent
	case COND:
		return OperandIdent
	case DATS:
		return OperandString
	case IASM:
		return OperandString
	case ADDB:
		return OperandNone
	case ADDW:
		return OperandNone
	case SUBB:
		return OperandNone
	case SUBW:
		return OperandNone
	case ANDB:
		return OperandNone
	case ANDW:
		return OperandNone
	case BORB:
		return OperandNone
	case BORW:
		return OperandNone
	case XORB:
		return OperandNone
	case XORW:
		return OperandNone
	case SHLB:
		return OperandNone
	case SHLW:
		return OperandNone
	case SHRB:
		return OperandNone
	case SHRW:
		return OperandNone
	default:
		return OperandNone
	}
}

// Instructions consist of an opcode and a single optional operand.
// This means some instructions have to look back at the previous
// instruction to get all their operands if one is not enough.
type Instruction struct {
	Operation // Operation

	Byte  uint8  // Filled in for OperandByte
	Word  uint16 // Filled in for OperandWord
	Int   int    // Filled in for OperandInt
	Str   string // Filled in for OperandString
	Ident string // Filled in for OperandIdent
}

// Program is a list of instructions
type Program []Instruction
