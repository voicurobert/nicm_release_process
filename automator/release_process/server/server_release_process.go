package server

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
)

type serverReleaseProcess struct {
	commands []commands.CommandInterface
}

var (
	ReleaseProcess = &serverReleaseProcess{}
)

func (c *serverReleaseProcess) Execute() error {
	for _, cmd := range c.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *serverReleaseProcess) Init() {

}

func (c *serverReleaseProcess) PrintCommands() {
	for _, c := range c.commands {
		c.Print()
	}
}

func (c *serverReleaseProcess) PrintOptions() {

}
