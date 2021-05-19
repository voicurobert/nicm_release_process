package app

import (
	"bufio"
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/process"
	"os"
)

func StartApplication() {
	fmt.Println("Started NICM Release Automator!")
	scanner := bufio.NewScanner(os.Stdin)
	startProcess := process.Process
	for {
		fmt.Printf("Command [%s]:", startProcess.GetName())
		scanner.Scan()
		textCommand := scanner.Text()

		switch textCommand {
		case process.HelpCommandText:
			startProcess.PrintPossibleProcesses()
		case
			process.ExitCommandText:
			closeApplication()
		case process.PrepareReleaseProcess,
			process.ClientReleaseProcess,
			process.ServerReleaseProcess,
			process.ActivateReleaseProcess,
			process.NIGReleaseProcess,
			process.GearsReleaseProcess,
			process.PreviousProcess:
			startProcess = startProcess.NextProcess(textCommand)
		case process.ExecuteProcess:
			startProcess = startProcess.Execute()
		case process.PrintCommands:
			startProcess.PrintCommands()
		case process.PrintOptions:
			startProcess.PrintOptions()
		default:
			fmt.Println("Unknown command...try again")
		}
	}
}

func closeApplication() {
	os.Exit(1)
}
