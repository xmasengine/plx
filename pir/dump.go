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

	switch i.Operation.Operand() {
	case OperandNone:
		fmt.Fprintln(wr)
	case OperandByte:
		fmt.Fprintln(wr, " ", i.Byte)
	case OperandWord:
		fmt.Fprintln(wr, " ", i.Word)
	case OperandInt:
		fmt.Fprintln(wr, " ", i.Int)
	case OperandIdent:
		fmt.Fprintln(wr, " ", i.Ident)

	case OperandString:
		fmt.Fprintln(wr, " ", strconv.QuoteToASCII(i.Str))
	default:
		return fmt.Errorf("unknown operand type for %s", i.String())
	}

	return nil
}
