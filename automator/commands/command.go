package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/voicurobert/nicm_release_process/automator/config"
	"strings"
)

type CommandInterface interface {
	Execute() error
	Print(int)
}

type command struct {
	Name     string
	Function func(...interface{}) error
	Params   []interface{}
}

func (c *command) Execute() error {
	if c.Function != nil {
		color.Yellow("executing command -> %s\n", c.Name)
		err := c.Function(c.Params...)
		if err != nil {
			color.Red("error when executing command: %s\n", err.Error())
			return nil
		}
		color.Green("Success!\n")
		color.Green(strings.Repeat("-#", 50))
		fmt.Println()
	}
	return nil
}

func NewCommandsFromConfig(cfgName string) []string {
	list := make([]string, 0)
	cfgMap := config.GetConfig()
	clientMap, ok := cfgMap[cfgName]
	if !ok {
		return list
	}
	for _, cmd := range clientMap {
		list = append(list, cmd)
	}
	return list
}

func (c *command) Print(tabs int) {
	tabChars := strings.Repeat(" ", tabs)
	color.Yellow("%s<!%s!> \n", tabChars, c.Name)
}

func (c *command) SetParams(params ...interface{}) {
	c.Params = params
}

func New(name string, function func(...interface{}) error, args ...interface{}) CommandInterface {
	return &command{
		Name:     name,
		Function: function,
		Params:   args,
	}
}
