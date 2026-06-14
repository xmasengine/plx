package pir

import "text/scanner"
import "io"
import "fmt"
import "errors"

type Position = scanner.Position

func Parse(rd io.Reader) (Program, error) {
	var prog Program
	var errs []error
	scan := &scanner.Scanner{}
	scan.Error = func(s *scanner.Scanner, msg string) {
		errs = append(errs, fmt.Errorf("%s:%s", s.Position, msg))
	}
	scan.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	scan.Mode = scanner.ScanIdents | scanner.ScanChars | scanner.ScanStrings |
		scanner.ScanRawStrings | scanner.ScanComments | scanner.SkipComments

	var ins Instruction
	err := ins.scan(scan)
	if err != nil {
		errs = append(errs, err)
	} else {
		prog = append(prog, ins)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return prog, nil
}

func (ins *Instruction) scan(scan *scanner.Scanner) error {
	tok := scan.Next()
	// scan identifier with name of PIR instruction
	if tok != scanner.Ident {
		return fmt.Errorf("%s:%s:%s", scan.Position, "unexpected token", scanner.TokenString(tok))
	}
	text := scan.TokenText()
	err := ins.Operation.UnmarshalText([]byte(text))
	if err != nil {
		return fmt.Errorf("%s:%w", scan.Position, err)
	}

	return nil
}
