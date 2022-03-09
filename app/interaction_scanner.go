package app

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/voicurobert/nicm_release_process/automator/interations"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"os"
	"strings"
)

var (
	historyPath        = make([]string, 1)
	currentInteraction = interations.MainInteraction
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

		if text == interations.HelpCommandText {
			handleHelpCommand(spaces)
			continue
		}
		if text == interations.ExitCommandText {
			os.Exit(1)
		}
		if strings.Contains(text, "release") || text == options.SetOptions {
			handleNewCommand(text)
			continue
		}
		if text == interations.PreviousProcess {
			handlePreviousCommand(text)
			continue
		}
		if text == interations.ExecuteProcess {
			handleExecuteCommand()
			continue
		}
		if text == interations.PrintCommands {
			currentInteraction.PrintCommands(spaces)
			continue
		}
		if text == options.PrintOptions {
			currentInteraction.PrintOptions(spaces)
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
		handlePreviousCommand(interations.PreviousProcess)
	}
}

func handleHelpCommand(spaces int) {
	currentInteraction.PrintPossibleInteraction(spaces)
}

func handleDefaultCommand(text string) {
	ok, err := setOption(text)
	if err != nil {
		color.Red(err.Error())
		return
	}
	if !ok {
		color.Red("unknown command...")
	}
}

func setOption(text string) (bool, error) {
	if currentInteraction.GetName() == options.SetOptions {
		if strings.Contains(text, options.SetOptionSeparator) {
			vec := strings.Split(text, options.SetOptionSeparator)
			return true, handleSetOptionsMethod(strings.TrimSpace(vec[0]), strings.TrimSpace(vec[1]))
		} else {
			return true, errors.New("separator not allowed, try using '=")
		}
	}
	return false, nil
}

func handleSetOptionsMethod(option, value string) error {
	if value == "" {
		return errors.New("empty value")
	}
	if option == options.SetWorkingPath {
		currentInteraction.SetWorkingPath(value)
		return nil
	}
	if option == options.SetGitPath {
		currentInteraction.SetGitPath(value)
		return nil
	}
	return errors.New(fmt.Sprintf("unknown option: %s", option))
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
