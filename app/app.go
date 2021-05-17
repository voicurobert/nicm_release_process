package app

import (
	"bufio"
	"fmt"
	"github.com/voicurobert/nicm_release_process/command"
	"os"
	"strings"
)

func StartApplication() {
	fmt.Println("Starting NICM Release Process Automator...Waiting for command... (use help to see possible commands)")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		textCommand := scanner.Text()
		switch textCommand {
		case command.HelpCommandText:
			command.ShowPossibleCommands()
		case command.ExitCommandText:
			closeApplication()
		default:
			fmt.Println("Unknown command...try again")
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}

func closeApplication() {
	os.Exit(1)
}
