package pir

import "fmt"
import "strconv"
import "io"
import "os"

func (p Program) DumpFile(name string) error {
	in, err := os.Create(name)
	if err != nil {
		return err
	}
	defer in.Close()
	return p.Dump(in)
}

func (p Program) Dump(wr io.Writer) error {
	for _, ins := range p {
		ins.Dump(wr)
	}
	return nil
}

func (i Instruction) Dump(wr io.Writer) error {
	fmt.Fprint(wr, i.Operation)

	ops := i.Operands

	for _, op := range ops {
		switch op.Kind {
		case KindNone:
			continue
		case KindByte:
			fmt.Fprint(wr, " ", op.Byte)
		case KindWord:
			fmt.Fprint(wr, " ", op.Word)
		case KindInt:
			fmt.Fprint(wr, " ", op.Int)
		case KindIdent:
			fmt.Fprint(wr, " ", op.Ident)
		case KindRegister:
			fmt.Fprint(wr, " ", op.Register)
		case KindTemporary:
			fmt.Fprint(wr, " ", op.Temporary)
		case KindString:
			fmt.Fprint(wr, " ", strconv.QuoteToASCII(op.Str))
		default:
			return fmt.Errorf("unknown operand type %d for %s", op.Kind, i.String())
		}
	}
	fmt.Fprintln(wr)

	return nil
}
