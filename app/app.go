package app

import (
	"bufio"
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/process"
	"github.com/voicurobert/nicm_release_process/automator/process2"
	"os"
)

func StartApplication() {
	fmt.Println("Started NICM Release Automator!")
	scanner := bufio.NewScanner(os.Stdin)
	startProcess := process.Process
	for {
		fmt.Printf("Waiting for command... (use help to see possible commands): %s-> \n", startProcess.Name)
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

func StartApp() {
	fmt.Println("Started NICM Release Automator!")
	scanner := bufio.NewScanner(os.Stdin)
	process2.Process.Run(scanner)
}

func closeApplication() {
	os.Exit(1)
}
