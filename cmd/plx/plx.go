// pxl is the driver of the PLX compiler.
// It can also call the several assemblers and emulators.
package main

import (
	"github.com/xmasengine/plx/arch"
	"github.com/xmasengine/plx/cli"
	"github.com/xmasengine/plx/plat"
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

func main() {
	c := common{CLI: cli.New()}
	c.registerFlags()
	err := c.CLI.Start()
	cli.ExitIfErr("", err)
}
