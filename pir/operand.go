//go:generate go tool go-enum --marshal
package pir

/*
Kind is a kind of operand.

ENUM(
None	// No operand, skipped.
Byte
Word
Int
Ident
String
Register
Temporary
Condition
)
*/
type Kind int

type Operand struct {
	Kind
	Byte      uint8     // Filled in for KindByte
	Word      uint16    // Filled in for KindWord
	Int       int       // Filled in for KindInt
	Str       string    // Filled in for KindString
	Ident     string    // Filled in for KindIdent
	Register  Register  // Filled in for KindRegister
	Temporary Temporary // Filled in for KindTemporary
	Condition Condition // Filled in for KindCondition
}
