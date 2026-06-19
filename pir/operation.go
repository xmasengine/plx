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

	R1L // R1L is register R1 low half byte register
	R1H // R1H is register R1 high half byte register
	R2L // R2L is register R2 low half byte register
	R2H // R2H is register R2 high half byte register
	R3L // R3L is register R3 low half byte register
	R3H // R3H is register R3 high half byte register
	R4L // R4L is register R4 low half byte register
	R4H // R4H is register R4 high half byte register

)
*/
type Half int

func (r Register) Halves() (lo, hi Half) {
	switch r {
	case R1:
		return R1L, R1H
	case R2:
		return R2L, R2H
	case R3:
		return R3L, R3H
	case R4:
		return R4L, R4H
	default:
		panic("unknown register")
	}
}

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

func (o Operation) Kinds() []Kind {
	switch o {
	case NOOP:
		return []Kind{}
	case MOVB:
		return []Kind{KindHalf, KindHalf}
	case MOVW:
		return []Kind{KindRegister, KindRegister}
	case INCB:
		return []Kind{KindHalf}
	case INCW:
		return []Kind{KindRegister}
	case DECB:
		return []Kind{KindHalf}
	case DECW:
		return []Kind{KindRegister}
	case PSHW:
		return []Kind{KindRegister}
	case POPW:
		return []Kind{KindRegister}
	case LITB:
		return []Kind{KindByte, KindHalf}
	case LITW:
		return []Kind{KindWord, KindHalf}
	case OUTB:
		return []Kind{KindHalf, KindInt}
	case OUTW:
		return []Kind{KindRegister, KindInt}
	case OUTA:
		return []Kind{KindRegister, KindHalf, KindInt}
	case INPB:
		return []Kind{KindInt, KindHalf}
	case INPW:
		return []Kind{KindInt, KindRegister}
	case INPA:
		return []Kind{KindInt, KindHalf, KindRegister}
	case ADDB:
		return []Kind{KindHalf, KindHalf}
	case ADDW:
		return []Kind{KindRegister, KindRegister}
	case SUBB:
		return []Kind{KindHalf, KindHalf}
	case SUBW:
		return []Kind{KindRegister, KindRegister}
	case ANDB:
		return []Kind{KindHalf, KindHalf}
	case ANDW:
		return []Kind{KindRegister, KindRegister}
	case BORB:
		return []Kind{KindHalf, KindHalf}
	case BORW:
		return []Kind{KindRegister, KindRegister}
	case XORB:
		return []Kind{KindHalf, KindHalf}
	case XORW:
		return []Kind{KindRegister, KindRegister}
	case SHLB:
		return []Kind{KindInt, KindHalf}
	case SHLW:
		return []Kind{KindInt, KindRegister}
	case SHRB:
		return []Kind{KindInt, KindHalf}
	case SHRW:
		return []Kind{KindInt, KindRegister}
	case LAND:
		return []Kind{KindIdent}
	case JUMP:
		return []Kind{KindIdent}
	case JPIF:
		return []Kind{KindCondition, KindIdent}
	case CMPB:
		return []Kind{KindCondition, KindHalf}
	case DATS:
		return []Kind{KindIdent, KindString}
	case IASM:
		return []Kind{KindString}
	case VARA:
		return []Kind{KindIdent, KindInt}
	case GEAB:
		return []Kind{KindRegister, KindHalf}
	case GEAW:
		return []Kind{KindRegister, KindRegister}
	case STAB:
		return []Kind{KindHalf, KindRegister}
	case STAW:
		return []Kind{KindRegister, KindRegister}
	case STOB:
		return []Kind{KindHalf, KindIdent}
	case STOW:
		return []Kind{KindRegister, KindIdent}
	case STOT:
		return []Kind{KindRegister, KindTemporary}
	case LOAD:
		return []Kind{KindIdent, KindRegister}
	case LOAB:
		return []Kind{KindIdent, KindHalf}
	case LOAW:
		return []Kind{KindIdent, KindRegister}
	case GETB:
		return []Kind{KindIdent, KindHalf}
	case GETW:
		return []Kind{KindIdent, KindRegister}
	case GETA:
		return []Kind{KindIdent, KindRegister}
	case GETT:
		return []Kind{KindTemporary, KindRegister}
	case FUNC:
		return []Kind{KindIdent}
	case CALL:
		return []Kind{KindIdent}
	case RETU:
		return []Kind{}
	case RETI:
		return []Kind{}
	case RETN:
		return []Kind{}
	case SINT:
		return []Kind{KindIdent}
	case SNMI:
		return []Kind{KindIdent}
	case COPY:
		return []Kind{KindRegister, KindHalf, KindRegister}
	case BANK:
		return []Kind{KindInt}
	case BATT:
		return []Kind{KindIdent}
	default:
		return []Kind{}
	}

}
