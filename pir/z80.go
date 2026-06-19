package pir

import "fmt"

// import "strconv"
import "io"
import "os"
import "bytes"
import "github.com/xmasengine/plx/plat"

/*

PIR -> z80 translator conventions:

We use the hardware stack as the virtual stack.

The register pairs HL, DE are used as the top (TOS) of stack and next
of stack (NXT), and over if these are used we spill to the stack.

The calling convetion for funtions and operator is the same,
they use the stack, but HL and DE are used for the first two parameters,
because the are TOS and NXT.

The A and F registers will be used by all operators so are always clobbered.
The BC register is free for other uses.

*/

func ParseFileZ80(name string, pl plat.Platform) (*bytes.Buffer, error) {
	pir, err := ParseFilePIR(name)
	if err != nil {
		return nil, err
	}
	return pir.EmitZ80(pl)
}

func (p Program) EmitZ80(pl plat.Platform) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	z80 := Z80{Emitter: Emitter{out: buf}, Platform: pl}
	err := z80.EmitProgram(p)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ParseFilesZ80(pl plat.Platform, output string, names ...string) error {
	prog := Program{}
	for _, name := range names {
		sub, err := ParseFilePIR(name)
		if err != nil {
			return err
		}
		prog = append(prog, sub...)
	}
	buf, err := prog.EmitZ80(pl)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, buf)
	return nil
}

type Emitter struct {
	out io.Writer
}

func (e *Emitter) Emitf(form string, args ...any) error {
	_, err := fmt.Fprintf(e.out, form+"\n", args...)
	return err
}

type Z80 struct {
	Emitter
	plat.Platform
	lastName  string
	lastLabel int
}

const z80Header = `
// Banking setup
banksize %#x
bankat 0
bank 0


// Boot section
org 0x0000
    di              //  disable interrupts
    im 1            //  Interrupt mode 1
    jp main         //  jump to main program

// Interrupt handler
org %#x
	reti // do nothing for now


//  Pause button/NMI handler
org %#x
    //  Do nothing
    retn

// Main program
main:
	// Set up stack.
    ld sp, %#x
`

func (z *Z80) emitHeader() error {
	info := z.Platform.Info()
	return z.Emitf(z80Header, info.BankSize, info.INT, info.NMI, info.Stack)
}

func (z *Z80) EmitProgram(p Program) error {
	z.emitHeader()
	for _, ins := range p {
		err := z.emitInstruction(ins)
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *Z80) emitInstruction(i Instruction) error {
	switch i.Operation {
	// no operation
	case NOOP:
		z.emitNOOP(i)
	// move register byte RnL -> RmL or RnH -> RmH
	case MOVB:
		z.emitMOVB(i)
	// move register word RnL -> RmL or RnH -> RmH
	case MOVW:
		z.emitMOVW(i)
	// increment register byte
	case INCB:
		z.emitINCB(i)
	// increment register word
	case INCW:
		z.emitINCW(i)
	// decrement register byte
	case DECB:
		z.emitDECB(i)
	// decrement register word
	case DECW:
		z.emitDECW(i)
	// push register word to data stack
	case PSHW:
		z.emitPSHW(i)
	// pop data stack word to register (there is no byte variant)
	case POPW:
		z.emitPOPW(i)
	// Store literal byte in register RnL or RnH
	case LITB:
		z.emitLITB(i)
	// Store literal word in register RnL
	case LITW:
		z.emitLITW(i)
	// Output register byte to port [int] (constant literal).
	case OUTB:
		z.emitOUTB(i)
	// Output register word to port [int] (constant literal).
	case OUTW:
		z.emitOUTW(i)
	// Output R1 must have the address, R2 the length, output to port [int]
	case OUTA:
		z.emitOUTA(i)
	// Input byte from port [int] to register.
	case INPB:
		z.emitINPB(i)
	// Input word from port [int] to register.
	case INPW:
		z.emitINPW(i)
	// Output R1 must have the address, R2 the length, input from port [int]
	case INPA:
		z.emitINPA(i)
	// Add byte Rn to Rm and store in Rm
	case ADDB:
		z.emitADDB(i)
	// Add word Rn to Rm and store in Rm
	case ADDW:
		z.emitADDW(i)
	// Subtract byte Rn from Rm and store in Rm
	case SUBB:
		z.emitSUBB(i)
	// Subtract word Rn from Rm and store in Rm
	case SUBW:
		z.emitSUBW(i)
	// AND byte Rn with Rm and store in Rm
	case ANDB:
		z.emitANDB(i)
	// AND word Rn with Rm and store in Rm
	case ANDW:
		z.emitANDW(i)
	// Binary OR byte Rn with Rm and store in Rm
	case BORB:
		z.emitBORB(i)
	// Binary OR word Rn with Rm and store in Rm
	case BORW:
		z.emitBORW(i)
	// Binary XOR byte Rn with Rm and store in Rm
	case XORB:
		z.emitXORB(i)
	// Binary XOR word Rn with Rm and store in Rm
	case XORW:
		z.emitXORW(i)
	// Shift left by Int to byte register.
	case SHLB:
		z.emitSHLB(i)
	// Shift left by Int to word register.
	case SHLW:
		z.emitSHLW(i)
	// Shift right by Int to byte register.
	case SHRB:
		z.emitSHRB(i)
	// Shift right by Int to word register.
	case SHRW:
		z.emitSHRW(i)
	// Define jump location where the jump may "land" [ident].
	case LAND:
		z.emitLAND(i)
	// Jump to tag [ident] unconditionally.
	case JUMP:
		z.emitJUMP(i)
	// Jump on [cond] to tag [ident] it TOS is TRUE, pop stack.
	case JPIF:
		z.emitJPIF(i)
	// Ident is one of [eq, gt, lt, etc], compare register with R1L and store in R1L.
	case CMPB:
		z.emitCMPB(i)
	// Data String.
	case DATS:
		z.emitDATS(i)
	// Inline assembly string.
	case IASM:
		z.emitIASM(i)
	// Allocate variable with name and size.
	case VARA:
		z.emitVARA(i)
	// Get byte indirectly from address pointed by Rn to Rm.
	case GEAB:
		z.emitGEAB(i)
	// Get word indirectly from address pointed by Rn to Rm.
	case GEAW:
		z.emitGEAW(i)
	// Store byte indirectly from address pointed by Rn from Rm.
	case STAB:
		z.emitSTAB(i)
	// Store word indirectly from address pointed by Rn from Rm.
	case STAW:
		z.emitSTAW(i)
	// Store byte in variable.
	case STOB:
		z.emitSTOB(i)
	// Store word in variable.
	case STOW:
		z.emitSTOW(i)
	// Store register to temporary Rn to Tn. There is no byte variant.
	case STOT:
		z.emitSTOT(i)
	// Get named data address to register.
	case LOAD:
		z.emitLOAD(i)
	// Get named data byte to register.
	case LOAB:
		z.emitLOAB(i)
	// Get named data word to register.
	case LOAW:
		z.emitLOAW(i)
	// Get named variable byte to register.
	case GETB:
		z.emitGETB(i)
	// Get named variable word to register.
	case GETW:
		z.emitGETW(i)
	// Get named variable address to register.
	case GETA:
		z.emitGETA(i)
	// Get temporary to register. There is no byte variant.
	case GETT:
		z.emitGETT(i)
	// Define callable function/sub.
	case FUNC:
		z.emitFUNC(i)
	// Call tag.
	case CALL:
		z.emitCALL(i)
	//  normally from a call.
	case RETU:
		z.emitRETU(i)
	//  from an interrupt call.
	case RETI:
		z.emitRETI(i)
	//  from an nmi call.
	case RETN:
		z.emitRETN(i)
	// Set tag as an interrupt handler.
	case SINT:
		z.emitSINT(i)
	// Set tag an an NMI handler.
	case SNMI:
		z.emitSNMI(i)
	// Copy R2 length bytes to R1. XXX: better ideas.
	case COPY:
		z.emitCOPY(i)
	// Switch active memory bank to constant int.
	case BANK:
		z.emitBANK(i)
	// Switch battery backed memory on or off and store address in R.
	case BATT:
		z.emitBATT(i)
	default:
		return fmt.Errorf("unknown instruction %d", i.Operation)
	}
	return nil
}

func (z Z80) pre(ident string) string {
	return "pir_" + ident
}

func (z *Z80) label() string {
	z.lastLabel++
	return fmt.Sprintf("pl_%d", z.lastLabel)
}

func (z *Z80) push() error {
	z.Emitf("push de")
	z.Emitf("ex de, hl")
	return nil
}

func (z *Z80) pop() error {
	z.Emitf("ex de,hl")
	z.Emitf("pop de")
	return nil
}

// no operation
func (z *Z80) emitNOOP(i Instruction) {}

// move register byte RnL -> RmL or RnH -> RmH
func (z *Z80) emitMOVB(i Instruction) {}

// move register word RnL -> RmL or RnH -> RmH
func (z *Z80) emitMOVW(i Instruction) {}

// increment register byte
func (z *Z80) emitINCB(i Instruction) {}

// increment register word
func (z *Z80) emitINCW(i Instruction) {}

// decrement register byte
func (z *Z80) emitDECB(i Instruction) {}

// decrement register word
func (z *Z80) emitDECW(i Instruction) {}

// push register word to data stack
func (z *Z80) emitPSHW(i Instruction) {}

// pop data stack word to register (there is no byte variant)
func (z *Z80) emitPOPW(i Instruction) {}

// Store literal byte in register RnL or RnH
func (z *Z80) emitLITB(i Instruction) {}

// Store literal word in register RnL
func (z *Z80) emitLITW(i Instruction) {}

// Output register byte to port [int] (constant literal).
func (z *Z80) emitOUTB(i Instruction) {}

// Output register word to port [int] (constant literal).
func (z *Z80) emitOUTW(i Instruction) {}

// Output R1 must have the address, R2 the length, output to port [int]
func (z *Z80) emitOUTA(i Instruction) {}

// Input byte from port [int] to register.
func (z *Z80) emitINPB(i Instruction) {}

// Input word from port [int] to register.
func (z *Z80) emitINPW(i Instruction) {}

// Output R1 must have the address, R2 the length, input from port [int]
func (z *Z80) emitINPA(i Instruction) {}

// Add byte Rn to Rm and store in Rm
func (z *Z80) emitADDB(i Instruction) {}

// Add word Rn to Rm and store in Rm
func (z *Z80) emitADDW(i Instruction) {}

// Subtract byte Rn from Rm and store in Rm
func (z *Z80) emitSUBB(i Instruction) {}

// Subtract word Rn from Rm and store in Rm
func (z *Z80) emitSUBW(i Instruction) {}

// AND byte Rn with Rm and store in Rm
func (z *Z80) emitANDB(i Instruction) {}

// AND word Rn with Rm and store in Rm
func (z *Z80) emitANDW(i Instruction) {}

// Binary OR byte Rn with Rm and store in Rm
func (z *Z80) emitBORB(i Instruction) {}

// Binary OR word Rn with Rm and store in Rm
func (z *Z80) emitBORW(i Instruction) {}

// Binary XOR byte Rn with Rm and store in Rm
func (z *Z80) emitXORB(i Instruction) {}

// Binary XOR word Rn with Rm and store in Rm
func (z *Z80) emitXORW(i Instruction) {}

// Shift left by Int to byte register.
func (z *Z80) emitSHLB(i Instruction) {}

// Shift left by Int to word register.
func (z *Z80) emitSHLW(i Instruction) {}

// Shift right by Int to byte register.
func (z *Z80) emitSHRB(i Instruction) {}

// Shift right by Int to word register.
func (z *Z80) emitSHRW(i Instruction) {}

// Define jump location where the jump may "land" [ident].
func (z *Z80) emitLAND(i Instruction) {}

// Jump to tag [ident] unconditionally.
func (z *Z80) emitJUMP(i Instruction) {}

// Jump on [cond] to tag [ident] it TOS is TRUE, pop stack.
func (z *Z80) emitJPIF(i Instruction) {}

// Ident is one of [eq, gt, lt, etc], compare register with R1L and store in R1L.
func (z *Z80) emitCMPB(i Instruction) {}

// Data String.
func (z *Z80) emitDATS(i Instruction) {}

// Inline assembly string.
func (z *Z80) emitIASM(i Instruction) {}

// Allocate variable with name and size.
func (z *Z80) emitVARA(i Instruction) {}

// Get byte indirectly from address pointed by Rn to Rm.
func (z *Z80) emitGEAB(i Instruction) {}

// Get word indirectly from address pointed by Rn to Rm.
func (z *Z80) emitGEAW(i Instruction) {}

// Store byte indirectly from address pointed by Rn from Rm.
func (z *Z80) emitSTAB(i Instruction) {}

// Store word indirectly from address pointed by Rn from Rm.
func (z *Z80) emitSTAW(i Instruction) {}

// Store byte in variable.
func (z *Z80) emitSTOB(i Instruction) {}

// Store word in variable.
func (z *Z80) emitSTOW(i Instruction) {}

// Store register to temporary Rn to Tn. There is no byte variant.
func (z *Z80) emitSTOT(i Instruction) {}

// Get named data address to register.
func (z *Z80) emitLOAD(i Instruction) {}

// Get named data byte to register.
func (z *Z80) emitLOAB(i Instruction) {}

// Get named data word to register.
func (z *Z80) emitLOAW(i Instruction) {}

// Get named variable byte to register.
func (z *Z80) emitGETB(i Instruction) {}

// Get named variable word to register.
func (z *Z80) emitGETW(i Instruction) {}

// Get named variable address to register.
func (z *Z80) emitGETA(i Instruction) {}

// Get temporary to register. There is no byte variant.
func (z *Z80) emitGETT(i Instruction) {}

// Define callable function/sub.
func (z *Z80) emitFUNC(i Instruction) {}

// Call tag.
func (z *Z80) emitCALL(i Instruction) {}

// Return normally from a call.
func (z *Z80) emitRETU(i Instruction) {}

// Return from an interrupt call.
func (z *Z80) emitRETI(i Instruction) {}

// Return from an nmi call.
func (z *Z80) emitRETN(i Instruction) {}

// Set tag as an interrupt handler.
func (z *Z80) emitSINT(i Instruction) {}

// Set tag an an NMI handler.
func (z *Z80) emitSNMI(i Instruction) {}

// Copy R2 length bytes to R1. XXX: better ideas.
func (z *Z80) emitCOPY(i Instruction) {}

// Switch active memory bank to constant int.
func (z *Z80) emitBANK(i Instruction) {}

// Switch battery backed memory on or off and store address in R.
func (z *Z80) emitBATT(i Instruction) {}
