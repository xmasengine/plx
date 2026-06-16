//go:generate go tool go-enum --marshal
package pir

/*
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
type Operand int
