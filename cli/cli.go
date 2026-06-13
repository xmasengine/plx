// package cli contains helpers for a command line with
// multiple sub commands
package cli

import "flag"
import "fmt"
import "os"
import "errors"
import "path/filepath"

// EnvDefault returns the named string like variable from env, or
// the default if not found.
func EnvDefault[S ~string](name string, def S) S {
	val, found := os.LookupEnv(name)
	if val == "" || !found {
		return def
	}
	return S(val)
}

// Fatalf prints and error message and exits.
func Fatalf(form string, args ...any) {
	fmt.Fprintf(os.Stderr, form+"\n", args...)
	os.Exit(3)
}

// ErrorExit prints multiple errors, like for a compiler output
// and exits. Does nothing if errors is empty.
func ErrorExit(msg string, errors ...error) {
	if len(errors) < 1 {
		return
	}
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(2)
}

// ExitIfErr does nothing if err i nil. If it is not nil, it prints the err and  if err is not nil
func ExitIfErr(msg string, err error) {
	if err == nil {
		return
	}
	if msg != "" {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, err.Error())
}

type Runner interface {
	Run(args ...string) error
}

type Command struct {
	*flag.FlagSet // we use Flagset.Name() as the name of the command also
	Runner        Runner
}

func (c *Command) Run(args ...string) error {
	err := c.FlagSet.Parse(args)
	if err != nil {
		return nil
	}
	return c.Runner.Run(c.FlagSet.Args()...)
}

func NewCommand(name string, runner Runner) *Command {
	return &Command{
		Runner:  runner,
		FlagSet: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

type CLI struct {
	*flag.FlagSet
	Commands []Command
	showHelp bool
	Default  *Command
}

func (c *CLI) Run(args ...string) error {
	if len(args) < 2 {
		if c.Default != nil {
			return c.Default.Run(args...)
		}
		c.help()
		return errors.New("please specify a sub command")
	}
	for _, com := range c.Commands {
		if com.Name() == args[1] {
			return com.Run(args...)
		}
	}
	return fmt.Errorf("unknown command %s", args[1])
}

func (c *CLI) PrintCommands() {
	for _, com := range c.Commands {
		fmt.Fprintf(os.Stderr, "\t%s\n", com.Name())
		com.FlagSet.PrintDefaults()
	}
}

func (c *CLI) helpCommand() *help {
	return &help{c: c}
}

func (c *CLI) help() error {
	fmt.Fprintf(os.Stderr, "Sub commands for %s:\n", filepath.Base(os.Args[0]))
	c.PrintCommands()
	fmt.Fprintf(os.Stderr, "Common flags for %s:\n", filepath.Base(os.Args[0]))
	c.FlagSet.PrintDefaults()
	return nil
}

type help struct {
	c *CLI
}

func (h *help) Run(args ...string) error {
	return h.c.help()
}

func (c *CLI) Command(name string, runner Runner) *Command {
	res := NewCommand(name, runner)
	c.Commands = append(c.Commands, *res)
	return res
}

func New() *CLI {
	res := &CLI{}
	res.FlagSet = flag.NewFlagSet("", flag.ContinueOnError)
	res.FlagSet.BoolVar(&res.showHelp, "h", false, "show help")
	res.Command("help", res.helpCommand())
	return res
}

// Start parses the command line arguments and runs the CLI.
func (c *CLI) Start() error {
	err := c.FlagSet.Parse(os.Args)
	if err != nil {
		return err
	}
	if c.showHelp {
		return c.help()
	}
	return c.Run(c.FlagSet.Args()...)
}
