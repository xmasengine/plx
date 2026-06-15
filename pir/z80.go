package pir

import "fmt"
import "strconv"
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
	case NOOP:
		return z.emitNOOP(i)
	case INCB:
		return z.emitINCB(i)
	case INCW:
		return z.emitINCW(i)
	case DECB:
		return z.emitDECB(i)
	case DECW:
		return z.emitDECW(i)
	case POPB:
		return z.emitPOPB(i)
	case POPW:
		return z.emitPOPW(i)
	case DUPB:
		return z.emitDUPB(i)
	case DUPW:
		return z.emitDUPW(i)
	case PSHB:
		return z.emitPSHB(i)
	case PSHW:
		return z.emitPSHW(i)
	case OUTB:
		return z.emitOUTB(i)
	case OUTW:
		return z.emitOUTW(i)
	case OUTA:
		return z.emitOUTA(i)
	case INPB:
		return z.emitINPB(i)
	case INPW:
		return z.emitINPW(i)
	case NAME:
		return z.emitNAME(i)
	case PSHT:
		return z.emitPSHT(i)
	case LABL:
		return z.emitLABL(i)
	case JUMP:
		return z.emitJUMP(i)
	case JPIF:
		return z.emitJPIF(i)
	case COND:
		return z.emitCOND(i)
	case DATS:
		return z.emitDATS(i)
	case IASM:
		return z.emitIASM(i)
	case ADDB:
		return z.emitADDB(i)
	case ADDW:
		return z.emitADDW(i)
	case SUBB:
		return z.emitSUBB(i)
	case SUBW:
		return z.emitSUBW(i)
	case ANDB:
		return z.emitANDB(i)
	case ANDW:
		return z.emitANDW(i)
	case BORB:
		return z.emitBORB(i)
	case BORW:
		return z.emitBORW(i)
	case XORB:
		return z.emitXORB(i)
	case XORW:
		return z.emitXORW(i)
	case SHLB:
		return z.emitSHLB(i)
	case SHLW:
		return z.emitSHLW(i)
	case SHRB:
		return z.emitSHRB(i)
	case SHRW:
		return z.emitSHRW(i)
	default:
		return fmt.Errorf("unknown instrction %d", i.Operation)
	}
	return nil
}

func (z Z80) pre(ident string) string {
	return "pir_" + ident
}

func (z *Z80) label() string {
	z.lastLabel++
	return fmt.Sprintf(".%d", z.lastLabel)
}

func (z *Z80) push() error {
	z.Emitf("push de")
	z.Emitf("ld de, hl")
	return nil
}

func (z *Z80) pop() error {
	z.Emitf("ld hl, de")
	z.Emitf("pop de")
	return nil
}

func (z *Z80) emitNOOP(i Instruction) error { return z.Emitf("nop") }

func (z *Z80) emitINCB(i Instruction) error { return z.Emitf("inc l") }

func (z *Z80) emitINCW(i Instruction) error { return z.Emitf("inc hl") }

func (z *Z80) emitDECB(i Instruction) error { return z.Emitf("dec l") }

func (z *Z80) emitDECW(i Instruction) error { return z.Emitf("inc hl") }

func (z *Z80) emitPOPB(i Instruction) error { return z.pop() }

func (z *Z80) emitPOPW(i Instruction) error { return z.pop() }

func (z *Z80) emitDUPB(i Instruction) error { return z.push() }

func (z *Z80) emitDUPW(i Instruction) error { return z.push() }

func (z *Z80) emitPSHB(i Instruction) error {
	z.push()
	z.Emitf("ld h,0")
	return z.Emitf("ld l,%d", i.Byte)
}

func (z *Z80) emitPSHW(i Instruction) error { z.push(); return z.Emitf("ld hl,%d", i.Word) }

func (z *Z80) emitOUTB(i Instruction) error {
	z.Emitf("ld a, l")
	z.Emitf("out (%x), l", i.Int)
	return z.pop()
}

func (z *Z80) emitOUTW(i Instruction) error {
	z.Emitf("ld a, l")
	z.Emitf("out (%x), a", i.Int)
	z.Emitf("ld a, h")
	z.Emitf("out (%x), l", i.Int)
	return z.pop()
}

func (z *Z80) emitOUTA(i Instruction) error {
	// LD should be the address, TOS, BC the count with the count in B.
	z.Emitf("ld c,%x", i.Int)
	z.Emitf("otir")
	z.pop()
	z.pop()
	return nil
}

func (z *Z80) emitINPB(i Instruction) error {
	z.push()
	z.Emitf("in l,(%x)", i.Int)
	return nil
}

func (z *Z80) emitINPW(i Instruction) error {
	z.push()
	z.Emitf("in l,(%x)", i.Int)
	z.Emitf("in h,(%x)", i.Int)
	return nil
}

func (z *Z80) emitNAME(i Instruction) error {
	z.Emitf("%s:", z.pre(i.Ident))
	return nil
}

func (z *Z80) emitPSHT(i Instruction) error {
	z.push()
	z.Emitf("ld hl,%s", z.pre(i.Ident))
	return nil
}

func (z *Z80) emitLABL(i Instruction) error {
	z.Emitf("%s:", z.pre(i.Ident))
	return nil
}

func (z *Z80) emitJUMP(i Instruction) error {
	z.Emitf("jmp %s", z.pre(i.Ident))
	return nil
}

func (z *Z80) emitJPIF(i Instruction) error {
	z.Emitf("ld a,l")
	z.pop()
	z.Emitf("cp a,0")
	z.Emitf("jp z, %s:", z.pre(i.Ident))
	return nil
}

func (z *Z80) emitCOND(i Instruction) error {
	trueLabel := z.label()
	endLabel := z.label()
	// todo: check conditions, 16 bits
	z.Emitf("ld a,l")
	z.pop()
	z.Emitf("cp a,l") // compare
	z.Emitf("jp %s,%s", z.pre(i.Ident), trueLabel)
	z.Emitf("ld hl,0") // false branch
	z.Emitf("jmp %s", endLabel)
	z.Emitf("%s\tld hl,1", trueLabel) // true branch
	z.Emitf("%s\tnop", endLabel)
	return nil
}

func (z *Z80) emitDATS(i Instruction) error {
	z.Emitf("ds %s", strconv.Quote(i.Str))
	return nil
}

func (z *Z80) emitIASM(i Instruction) error {
	z.Emitf("\n%s\n", i.Str)
	return nil
}

func (z *Z80) emitADDB(i Instruction) error {
	z.Emitf("ld a,l")
	z.pop()
	z.Emitf("add a,l")
	z.Emitf("ld l,a")
	z.Emitf("ld h,0")
	// no need to do anything else, l has now the result.
	return nil
}

func (z *Z80) emitADDW(i Instruction) error {
	z.Emitf("add hl,de") // TOS + NEXT -> TOS
	z.Emitf("pop de")    // get next of stack, drop de
	// no need to do anything else, hl has now the result.
	return nil
}

func (z *Z80) emitSUBB(i Instruction) error {
	return nil
	z.Emitf("ld a,l")
	z.pop()
	z.Emitf("sub a,l")
	z.Emitf("ld l,a")
	z.Emitf("ld h,0")
	// no need to do anything else, l has now the result.
	return nil
}

func (z *Z80) emitSUBW(i Instruction) error {
	z.Emitf("sbc hl,de") // TOS - NEXT -> TOS
	z.Emitf("pop de")    // get next of stack, drop de
	// no need to do anything else, hl has now the result.
	return nil
}

func (z *Z80) emitANDB(i Instruction) error {
	z.Emitf("ld a,l")
	z.pop()
	z.Emitf("and a,l")
	z.Emitf("ld l,a")
	z.Emitf("ld h,0")
	// no need to do anything else, l has now the result.
	return nil
}

func (z *Z80) emitANDW(i Instruction) error { return nil }

func (z *Z80) emitBORB(i Instruction) error { return nil }

func (z *Z80) emitBORW(i Instruction) error { return nil }

func (z *Z80) emitXORB(i Instruction) error { return nil }

func (z *Z80) emitXORW(i Instruction) error { return nil }

func (z *Z80) emitSHLB(i Instruction) error { return nil }

func (z *Z80) emitSHLW(i Instruction) error { return nil }

func (z *Z80) emitSHRB(i Instruction) error { return nil }

func (z *Z80) emitSHRW(i Instruction) error { return nil }
