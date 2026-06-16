//go:generate go tool go-enum --noprefix --marshal --names --nocomments
package pir

/*
ENUM(

	NOOP                  // no operation
	INCB                  // increment stop of stack byte
	INCW                  // increment stop of stack word
	DECB                  // increment stop of stack byte
	DECW                  // increment stop of stack word
	POPB                  // Drop byte from data stack
	POPW                  // Drop word from data stack
	DUPB                  // Duplicate byte on top of stack
	DUPW                  // Duplicate word on top of stack
	NXTW                  // Duplicate word on next of stack to top of stack
	PSHB // Push byte literal [byte] to data stack.
	PSHW // Push word literal [word] to data stack.
	PSHA // Push address of [ident] to data stack.
	OUTB // Output TOS byte to port [int] (constant literal) and pop.
	OUTW // Output TOS word to port [int] (constant literal) and pop.
	OUTA // Output TOS address, length NXT bytes to port [int] (constant literal) and pop twice.
	INPB // Input byte from port [int], push to stack.
	INPW // Input word from port [int], push to stack.
	NAME // Name tag of next DATS, VARI, FUNC, etc instruction [ident].
	PSHT // Push address of tag [ident] to data stack.
	LABL // Define jump location [ident].
	JUMP // Jump to tag [ident] unconditionally.
	JPIF // Jump to tag [ident] it TOS is TRUE, pop stack.
	COND // Ident is one of [eq, gt, lt, etc], compare and push boolean to TOS.
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
	GETB

)
*/
type Operation int
