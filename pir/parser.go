package pir

import "text/scanner"
import "io"
import "fmt"
import "errors"
import "strconv"
import "os"

type Position = scanner.Position

func ParseFilePIR(name string) (Program, error) {
	in, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	return Parse(in, name)
}

func ParseFilesPIR(output string, names ...string) error {
	prog := Program{}
	for _, name := range names {
		sub, err := ParseFilePIR(name)
		if err != nil {
			return err
		}
		prog = append(prog, sub...)
	}
	return prog.DumpFile(output)
}

func Parse(rd io.Reader, name string) (Program, error) {
	var prog Program
	var errs []error
	scan := &scanner.Scanner{}
	scan.Init(rd)
	scan.Position.Filename = name
	scan.Mode = scanner.ScanIdents |
		scanner.ScanInts |
		scanner.ScanChars |
		scanner.ScanStrings |
		scanner.ScanRawStrings |
		scanner.ScanComments |
		scanner.SkipComments

	scan.Error = func(s *scanner.Scanner, msg string) {
		errs = append(errs, fmt.Errorf("scanner error:%s:%s", s.Position, msg))
	}
	scan.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

	for {
		tok := scan.Peek()
		if tok == scanner.EOF {
			break
		}
		var ins Instruction
		err := ins.scan(scan)
		if err != nil {
			errs = append(errs, err)
		} else if ins.Operation != SKIP {
			prog = append(prog, ins)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return prog, nil
}

func scanError(scan *scanner.Scanner, msg string, tok rune) error {
	return fmt.Errorf("%s:%s:%s:%s", scan.Position, msg, scanner.TokenString(tok), scan.TokenText())
}

func accept(scan *scanner.Scanner, accept ...rune) (token rune, text string, err error) {
	tok := scan.Scan()
	for _, r := range accept {
		if tok == r {
			return tok, scan.TokenText(), nil
		}
	}
	exp := "expected"
	for _, r := range accept {
		exp += " " + scanner.TokenString(r)
	}
	err = scanError(scan, "unexpected token "+exp, tok)
	return 0, "", err
}

func (ins *Instruction) scan(scan *scanner.Scanner) error {
	tok, text, err := accept(scan, scanner.Ident, '\n')
	if err != nil {
		return err
	}
	if tok == '\n' {
		ins.Operation = SKIP
		return nil
	}

	err = ins.Operation.UnmarshalText([]byte(text))
	if err != nil {
		return scanError(scan, err.Error(), tok)
	}

	switch ins.Operation.Operand() {
	case OperandNone:
		break
	case OperandByte:
		tok, text, err := accept(scan, scanner.Char, scanner.Int)
		if err != nil {
			return err
		}
		if tok == scanner.Char {
			r, _, _, err := strconv.UnquoteChar(text[1:len(text)], '\'')
			if err != nil {
				return scanError(scan, err.Error(), tok)
			}
			ins.Byte = byte(r)
		} else {
			num, err := strconv.ParseUint(text, 0, 8)
			if err != nil {
				return scanError(scan, err.Error(), tok)
			}
			ins.Byte = uint8(num)
		}
	case OperandWord:
		tok, text, err := accept(scan, scanner.Int)
		if err != nil {
			return err
		}
		num, err := strconv.ParseUint(text, 0, 16)
		if err != nil {
			return scanError(scan, err.Error(), tok)
		}
		ins.Word = uint16(num)
	case OperandInt:
		tok, text, err := accept(scan, scanner.Int)
		if err != nil {
			return err
		}
		num, err := strconv.ParseInt(text, 0, 0)
		if err != nil {
			return scanError(scan, err.Error(), tok)
		}
		ins.Int = int(num)

	case OperandIdent:
		_, text, err := accept(scan, scanner.Ident)
		if err != nil {
			return err
		}
		ins.Ident = text

	case OperandString:
		tok, text, err := accept(scan, scanner.String, scanner.RawString)
		if err != nil {
			return err
		}
		str, err := strconv.Unquote(text)
		if err != nil {
			scanError(scan, err.Error(), tok)
		}
		ins.Str = str

	default:
		return fmt.Errorf("%s:%s", scan.Position, "unknown operand type")
	}
	_, _, err = accept(scan, ';', '\n')

	return err
}
