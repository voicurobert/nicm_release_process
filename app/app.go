package app

import (
	"bufio"
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/options"
	"github.com/voicurobert/nicm_release_process/automator/process"
	"os"
	"strings"
)

var (
	historyPath  = make([]string, 1)
	startProcess = process.Process
)

func StartApplication() {
	fmt.Println("Started NICM Release Automator!")
	scanner := bufio.NewScanner(os.Stdin)
	historyPath[0] = startProcess.GetName()
	for {
		fmt.Printf("Command [%s]:", getHistoryString())
		scanner.Scan()
		text := scanner.Text()
		if text == process.HelpCommandText {
			handleHelpCommand()
			continue
		}
		if text == process.ExitCommandText {
			handleHelpCommand()
			continue
		}
		if strings.Contains(text, "release") || strings.Contains(text, "set") {
			handleNewCommand(text)
			continue
		}
		if text == process.PreviousProcess {
			handlePreviousCommand(text)
			continue
		}
		if text == process.ExecuteProcess {
			handleExecuteCommand()
		}
		if text == process.PrintCommands {
			startProcess.PrintCommands()
			continue
		}
		if text == options.PrintOptions {
			startProcess.PrintOptions()
			continue
		}
		handleDefaultCommand(text)
	}
}

func handleNewCommand(text string) {
	startProcess = startProcess.NextProcess(text)
	addHistory(startProcess.GetName())
}

func handleDefaultCommand(text string) {
	prevProcess := setOption(startProcess, text)
	if prevProcess == nil {
		fmt.Println("Unknown command...")
	} else {
		startProcess = prevProcess
		historyPath = historyPath[:len(historyPath)-1]
	}
}
func handlePreviousCommand(text string) {
	prevProcess := startProcess.NextProcess(text)
	if prevProcess != startProcess {
		startProcess = prevProcess
		historyPath = historyPath[:len(historyPath)-1]
	}
}

func handleExecuteCommand() {
	startProcess = startProcess.Execute()
	addHistory(startProcess.GetName())
}

func handleHelpCommand() {
	startProcess.PrintPossibleProcesses()
}

func setOption(currentProcess process.UserProcessInterface, text string) process.UserProcessInterface {
	if currentProcess.GetName() == options.SetWorkingPath {
		currentProcess.SetWorkingPath(text)
		return currentProcess.NextProcess(process.PreviousProcess)
	}
	return nil
}

func addHistory(name string) {
	historyPath = append(historyPath, name)
}

func getHistoryString() string {
	var sb strings.Builder
	for idx, str := range historyPath {
		if idx != 0 {
			sb.WriteString("->")
		}
		sb.WriteString(str)
	}
	return sb.String()
}
