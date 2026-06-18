//go:generate go tool go-enum --noprefix --marshal --names --nocomments
package pir

/*
ENUM(

	NOOPE // no operation
	MOVRB // move register byte RnL -> RmL or RnH -> RmH
	MOVRW // move register word RnL -> RmL or RnH -> RmH
	GETRB // Get byte from address pointed by Rn to Rm
	GETRW // Get word from address pointed by Rn to Rm
	INCRB // increment register byte
	INCRW // increment register word
	DECRB // decrement register byte
	DECRW // decrement register word
	PUSHW // push register word to data stack
	POPRW // pop data stack word to register (there is no byte variant)
	LITRB // Store literal byte in register RnL or RnH
	LITRW // Store literal word in register RnL
	OUTRB // Output register byte to port [int] (constant literal).
	OUTRW // Output register word to port [int] (constant literal).
	OUTRA // Output R1 must have the address, R2 the length, output to port [int]
	INPRB // Input byte from port [int] to register.
	INPWB // Input word from port [int] to register.
	INPRA // Output R1 must have the address, R2 the length, input from port [int]
	NAMET // Name tag of next DATS, VARI, FUNC, etc instruction [ident].
	LOADT // Load tag address into register.
	LABEL // Define jump location [ident].
	JUMPT // Jump to tag [ident] unconditionally.
	JPIFT // Jump on [cond] to tag [ident] it TOS is TRUE, pop stack.
	CONDB // Ident is one of [eq, gt, lt, etc], compare register with R1L and store in R1L.
	DATAS // Data String.
	IASMS // Inline assembly string.
	ADDRB // Add byte Rn to Rm and store in Rm
	ADDRW // Add word Rn to Rm and store in Rm
	SUBRB // Subtract byte Rn from Rm and store in Rm
	SUBRW // Subtract word Rn from Rm and store in Rm
	ANDRB // AND byte Rn with Rm and store in Rm
	ANDRW // AND word Rn with Rm and store in Rm
	BORRB // Binary OR byte Rn with Rm and store in Rm
	BORRW // Binary OR word Rn with Rm and store in Rm
	XORRB // Binary XOR byte Rn with Rm and store in Rm
	XORRW // Binary XOR word Rn with Rm and store in Rm
	SHLRB // Shift left by Int to byte register.
	SHLRW // Shift left by Int to word register.
	SHRRB // Shift right by Int to byte register.
	SHRRW // Shift right by Int to word register.
	STORT // Store register to temporary Rn to Tn. There is no byte variant.
	GETTR // Get temporary to register. There is no byte variant.
	CALLT // Call tag.
	RETRN // Return normally from a call.
	RETIN // Return from an interrupt call.
	RETNM // Return from an nmi call.
	INTRT // Set tag as an interrupt handler.
	NMIHT // Set tag an an NMI handler.
	COPYR // Copy R2 length bytes to R1. XXX: better ideas.
	BANKI // Switch active memory bank to constant int.
	BATER // Switch battery backed memory on or off and store address in R.

)
*/
type Operation int
