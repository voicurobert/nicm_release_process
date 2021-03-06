package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/voicurobert/nicm_release_process/automator/config"
	"strconv"
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

	cfgMap := config.GetConfig()
	clientMap, ok := cfgMap[cfgName]
	list := make([]string, len(clientMap))

	if !ok {
		return list
	}
	for idx, cmd := range clientMap {
		intIdx, _ := strconv.Atoi(idx)
		list[intIdx-1] = cmd
		//list = append(list, cmd)
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
