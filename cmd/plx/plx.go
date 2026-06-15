// pxl is the driver of the PLX compiler.
// It can also call the several assemblers and emulators.
package main

import "errors"

import (
	"github.com/xmasengine/plx/arch"
	"github.com/xmasengine/plx/cli"
	"github.com/xmasengine/plx/pir"
	"github.com/xmasengine/plx/plat"
	"github.com/xmasengine/plx/z80asm"
)

type common struct {
	arch.Architecture // architecture, e.g. z80, 6052, etc.
	plat.Platform     // platform, e.g. NES, SMS, EMU
	*cli.CLI
}

func (c *common) registerFlags() {
	c.CLI.TextVar(&c.Architecture, "A", arch.Z80, "archictecture (z80, 6052)")
	c.CLI.TextVar(&c.Platform, "P", plat.SMS, "archictecture (sms, nes)")
}

type asm struct {
	cli.Command
	output string
	common *common
}

func (a *asm) Prepare(c *cli.Command) error {
	c.StringVar(&a.output, "o", "", "assembler output file name")
	return nil
}

func (a *asm) Run(args ...string) error {
	if a.output == "" {
		return errors.New("please specify assembler output file")
	}
	if len(args) < 1 {
		return errors.New("please specify assembler source files")
	}

	if a.common.Architecture == arch.Z80 {
		if a.common.Platform == plat.SMS {
			return z80asm.AssembleSMS(a.output, args)
		} else if a.common.Platform == plat.BIN {
			return z80asm.AssembleFiles(a.output, args)
		} else {
			return errors.New("TODO: platform not supported yet")
		}
	}
	return errors.New("TODO: architecture not supported yet")
}

type pirc struct {
	cli.Command
	output string
	common *common
}

func (p *pirc) Prepare(c *cli.Command) error {
	c.StringVar(&p.output, "o", "", "pir output file name")
	return nil
}

func (p *pirc) Run(args ...string) error {
	if p.output == "" {
		return errors.New("please specify pir output file")
	}
	if len(args) < 1 {
		return errors.New("please specify pir source files")
	}
	if p.common.Architecture == arch.PIR {
		return pir.ParseFilesPIR(p.output, args...)
	} else if p.common.Architecture == arch.Z80 {
		return pir.ParseFilesZ80(p.common.Platform, p.output, args...)
	} else {
		return errors.New("achitecture not yet supported: " + p.common.Architecture.String())
	}
}

func main() {
	c := &common{CLI: cli.New()}
	c.registerFlags()
	a := &asm{common: c}
	p := &pirc{common: c}
	c.CLI.Command("asm", a)
	c.CLI.Command("pir", p)
	err := c.CLI.Start()
	cli.ExitIfErr("", err)
}
