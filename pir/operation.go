//go:generate go tool go-enum --noprefix --marshal --names --nocomments
package pir

/*
ENUM(

	NOOP // no operation
	MOVB // move register byte RnL -> RmL or RnH -> RmH
	MOVW // move register word RnL -> RmL or RnH -> RmH
	INCB // increment register byte
	INCW // increment register word
	DECB // decrement register byte
	DECW // decrement register word
	PSHW // push register word to data stack
	POPW // pop data stack word to register (there is no byte variant)
	LITB // Store literal byte in register RnL or RnH
	LITW // Store literal word in register RnL
	OUTB // Output register byte to port [int] (constant literal).
	OUTW // Output register word to port [int] (constant literal).
	OUTA // Output R1 must have the address, R2 the length, output to port [int]
	INPB // Input byte from port [int] to register.
	INPW // Input word from port [int] to register.
	INPA // Output R1 must have the address, R2 the length, input from port [int]
	ADDB // Add byte Rn to Rm and store in Rm
	ADDW // Add word Rn to Rm and store in Rm
	SUBB // Subtract byte Rn from Rm and store in Rm
	SUBW // Subtract word Rn from Rm and store in Rm
	ANDB // AND byte Rn with Rm and store in Rm
	ANDW // AND word Rn with Rm and store in Rm
	BORB // Binary OR byte Rn with Rm and store in Rm
	BORW // Binary OR word Rn with Rm and store in Rm
	XORB // Binary XOR byte Rn with Rm and store in Rm
	XORW // Binary XOR word Rn with Rm and store in Rm
	SHLB // Shift left by Int to byte register.
	SHLW // Shift left by Int to word register.
	SHRB // Shift right by Int to byte register.
	SHRW // Shift right by Int to word register.
	LAND // Define jump location where the jump may "land" [ident].
	JUMP // Jump to tag [ident] unconditionally.
	JPIF // Jump on [cond] to tag [ident] it TOS is TRUE, pop stack.
	CMPB // Ident is one of [eq, gt, lt, etc], compare register with R1L and store in R1L.
	DATS // Data String.
	IASM // Inline assembly string.
	VARA // Allocate variable with name and size.
	GEAB // Get byte indirectly from address pointed by Rn to Rm.
	GEAW // Get word indirectly from address pointed by Rn to Rm.
	STAB // Store byte indirectly from address pointed by Rn from Rm.
	STAW // Store word indirectly from address pointed by Rn from Rm.
	STOB // Store byte in variable.
	STOW // Store word in variable.
	STOT // Store register to temporary Rn to Tn. There is no byte variant.
	LOAD // Get named data address to register.
	LOAB // Get named data byte to register.
	LOAW // Get named data word to register.
	GETB // Get named variable byte to register.
	GETW // Get named variable word to register.
	GETA // Get named variable address to register.
	GETT // Get temporary to register. There is no byte variant.
	FUNC // Define callable function/sub.
	CALL // Call tag.
	RETU // Return normally from a call.
	RETI // Return from an interrupt call.
	RETN // Return from an nmi call.
	SINT // Set tag as an interrupt handler.
	SNMI // Set tag an an NMI handler.
	COPY // Copy R2 length bytes to R1. XXX: better ideas.
	BANK // Switch active memory bank to constant int.
	BATT // Switch battery backed memory on or off and store address in R.

)
*/
type Operation int

func (o Operation) OperandKinds() []Kind {
	switch o {

	default:
		return []Kind{}
	}
}

/*
ENUM(

	R1 // R1 is the virtual accumulator register
	R2 // R2 is the virtual counter register
	R3 // R3 is the virtual pointer register
	R4 // R4 is the virtual backup register

)
*/
type Register int

/*
ENUM(

	T1 // T1 is one of four temporary virtual registers.
	T2 // T2 is one of four temporary virtual registers.
	T3 // R3 is one of four temporary virtual registers.
	T4 // R4 is one of four temporary virtual registers.

)
*/
type Temporary int

/*
ENUM(

	ZERO // Zero condition
	NOZE // Not zero condition
	LEST // Less than condition
	GRTT // Greater than condition
	LESE // Less than condition
	GETE // Greater than condition
	EQUA // Equal condition

)
*/
type Condition int
