package cli

import (
	"fmt"
	"strings"
)

type ParsedCmd struct {
	Cmd  []string
	Args []string
}

func (p *ParsedCmd) CmdString(conf *Conf) string {
	return strings.Join(p.Cmd, conf.Sep)
}

type Cmd struct {
	Name    string
	Func    func(...string) error
	SubCmds map[string]*Cmd
}

func (c *Cmd) Add(cmd *Cmd) {
	c.SubCmds[cmd.Name] = cmd
}

func CreateCmdRoot() *Cmd {
	return &Cmd{SubCmds: make(map[string]*Cmd)}
}
func CreateCmd(name string, f func(...string) error) *Cmd {
	return &Cmd{Name: name, Func: f, SubCmds: make(map[string]*Cmd)}
}
func CreateCmdNS(name string) *Cmd {
	return &Cmd{Name: name, SubCmds: make(map[string]*Cmd)}
}

func (c *Cmd) Get(arg []string) *Cmd {
	if c == nil {
		return nil
	} else if len(arg) == 0 {
		return c
	} else {
		cmd := c.SubCmds[arg[0]]
		return cmd.Get(arg[1:])
	}
}

func (c *Cmd) Call(parsed *ParsedCmd, conf *Conf) error {
	cmd := c.Get(parsed.Cmd)
	if cmd == nil {
		return fmt.Errorf("Command not found: %v", parsed.CmdString(conf))
	} else if cmd.Func == nil {
		if len(parsed.Cmd) == 0 {
			return fmt.Errorf("Missing command")
		} else {
			return fmt.Errorf("%v is not a valid command", parsed.CmdString(conf))
		}
	} else {
		return cmd.Func(parsed.Args...)
	}
}

var MainCmd = CreateCmdRoot()

type Conf struct {
	Sep string
}

var defaultConf Conf = Conf{Sep: ":"}

func prepareConf(conf *Conf) *Conf {
	newConf := &Conf{}
	if len(conf.Sep) == 0 {
		newConf.Sep = defaultConf.Sep
	} else {
		newConf.Sep = conf.Sep
	}
	return newConf
}

func Parse(args []string) error {
	return Parse2(args, &Conf{})
}

func Parse2(args []string, conf *Conf) error {
	conf = prepareConf(conf)
	cmd := &ParsedCmd{}
	if len(args) > 0 {
		cmd.Cmd = strings.Split(args[0], conf.Sep)
		cmd.Args = args[1:]
	}
	return MainCmd.Call(cmd, conf)
}
