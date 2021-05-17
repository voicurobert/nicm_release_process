package command

import "fmt"

const (
	HelpCommandText = "help"
	ExitCommandText = "exit"
)

var (
	CommandsList = []string{HelpCommandText, ExitCommandText}
)

func ShowPossibleCommands() {
	fmt.Println("Commands:")
	for _, command := range CommandsList {
		fmt.Printf("\t%s\n", command)
	}
}
