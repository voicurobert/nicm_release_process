package app

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/voicurobert/nicm_release_process/automator"
	"os"
	"strings"
)

var (
	historyPath        = make([]string, 1)
	currentInteraction = automator.MainInteraction
)

func StartInteracting() {
	scanner := bufio.NewScanner(os.Stdin)
	historyPath[0] = currentInteraction.GetName()
	for {
		color.Cyan("Command [%s]:", getHistoryString())
		scanner.Scan()
		text := scanner.Text()
		spaces := 9 + (len(historyPath) * 2)
		for _, path := range historyPath {
			spaces += len(path)
		}

		if strings.Contains(text, "release") {
			handleNewCommand(text)
			continue
		}

		if text == automator.HelpCommandText {
			handleHelpCommand(spaces)
			continue
		}

		if text == automator.ExitCommandText {
			os.Exit(1)
		}

		if text == automator.PreviousProcess {
			handlePreviousCommand(text)
			continue
		}

		if text == automator.ExecuteProcess {
			handleExecuteCommand()
			continue
		}

		if text == automator.PrintCommands {
			currentInteraction.PrintCommands(spaces)
			continue
		}

		handleDefaultCommand(text)
	}
}

func handleNewCommand(text string) {
	nextInter := currentInteraction.NextInteraction(text)
	if nextInter.GetName() == currentInteraction.GetName() {
		return
	}
	currentInteraction = nextInter
	addHistory(currentInteraction.GetName())
}

func handlePreviousCommand(text string) {
	prevProcess := currentInteraction.NextInteraction(text)
	if prevProcess != currentInteraction {
		currentInteraction = prevProcess
		historyPath = historyPath[:len(historyPath)-1]
	}
}

func handleExecuteCommand() {
	ok := currentInteraction.Execute()
	if ok == true {
		handlePreviousCommand(automator.PreviousProcess)
	}
}

func handleHelpCommand(spaces int) {
	currentInteraction.PrintPossibleInteraction(spaces)
}

func handleDefaultCommand(text string) {
	color.Red("unknown command %s...", text)
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
