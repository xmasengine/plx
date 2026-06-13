// pxl is the driver of the PLX compiler.
// It can also call the several assemblers and emulators.
package main

import "flag"
import "fmt"
import "os"

import (
	"github.com/xmasengine/plx/arch"
	"github.com/xmasengine/plx/plat"
)

func envDefault[S ~string](name string, def S) S {
	val, found := os.LookupEnv(name)
	if val == "" || !found {
		return def
	}
	return S(val)
}

func fatalf(form string, args ...any) {
	fmt.Fprintf(os.Stderr, form+"\n", args...)
	os.Exit(3)
}

func errorExit(msg string, errors ...error) {
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(2)
}

type common struct {
	arch.Architecture      // architecture, e.g. z80, 6052, etc.
	plat.Platform          // platform, e.g. NES, SMS, EMU
	help              bool // show help
}

func (c *common) registerFlags() {
	flag.TextVar(&c.Architecture, "A", arch.Z80, "archictecture (z80, 6052)")
	flag.TextVar(&c.Platform, "P", plat.SMS, "archictecture (sms, nes)")
}

type command struct {
	common
	name  string
	usage string
	set   flag.FlagSet
	run   func(args ...string)
}

type cli struct {
	common
	commands []command
}

func (c *cli) registerFlags() {
	c.common.registerFlags()
}

func (c *cli) run(args []string) {
	if len(args) < 1 {
		fatalf("please specify a sub command")
	}
	for _, com := range c.commands {
		if com.name == args[0] {
			com.run(args...)
			return
		}
	}
	fatalf("unknown command %s", args[0])
}

func (c *cli) printCommands() {
	for _, com := range c.commands {
		fmt.Fprintf(os.Stderr, "\t%s\t%s\n", com.name, com.usage)
		com.set.PrintDefaults()
	}
}

func (c *cli) help(args ...string) {
	fmt.Fprintf(os.Stderr, "Sub commands for %s:\n", os.Args[0])
	c.printCommands()
	fmt.Fprintf(os.Stderr, "Common Flags for %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func (c *cli) command(name, usage string, run func(args ...string)) *command {
	res := &command{
		name:  name,
		usage: usage,
		run:   run,
	}
	c.commands = append(c.commands, *res)
	return res
}

func newCli() *cli {
	res := &cli{}
	res.command("help", "show this help", res.help)
	return res
}

func main() {
	c := newCli()
	c.registerFlags()
	flag.Parse()
	if c.common.help {
		c.help()
	}
	c.run(flag.Args())
}
