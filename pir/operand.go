//go:generate go tool go-enum --marshal
package pir

/*
Kind is a kind of operand.

ENUM(
None
Byte
Word
Int
Ident
String
Register
Temporary
)
*/
type Kind int

type Operand struct {
	Kind
	Byte      uint8  // Filled in for KindByte
	Word      uint16 // Filled in for KindWord
	Int       int    // Filled in for KindInt
	Str       string // Filled in for KindString
	Ident     string // Filled in for KindIdent
	Register  string // Filled in for KindRegister
	Temporary string // Filled in for KindTemporary
}
