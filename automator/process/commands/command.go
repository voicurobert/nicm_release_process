package commands

import (
	"fmt"
	"strings"
)

type CommandInterface interface {
	Execute() error
	Print()
}

type Command struct {
	Name     string
	Function func(...interface{}) error
	Params   []interface{}
}

func (command *Command) Execute() error {
	if command.Function != nil {
		fmt.Printf("executing command -> %s\n", command.Name)
		err := command.Function(command.Params...)
		if err != nil {
			fmt.Printf("error when executing command: %s\n", err.Error())
			return nil
		}
		fmt.Println("success!")
		fmt.Println(strings.Repeat("-#", 50))
	}
	return nil
}

func (command *Command) Print() {
	fmt.Printf("\t <!%s!> \n", command.Name)
}

func (command *Command) SetParams(params ...interface{}) {
	command.Params = params
}

func NewCommand(name string, function func(...interface{}) error, args ...interface{}) CommandInterface {
	return &Command{
		Name:     name,
		Function: function,
		Params:   args,
	}
}
