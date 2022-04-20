package automator

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const (
	HelpCommandText      = "help"
	ExitCommandText      = "exit"
	NICMProcess          = "nicm_process"
	ClientReleaseProcess = "client_release"
	ServerReleaseProcess = "server_release"
	NIGReleaseProcess    = "nig_release"
	GearsReleaseProcess  = "gears_release"
	DTSTestRelease       = "dts_test_release"
	DTSProdRelease       = "dts_prod_release"
	PreviousProcess      = "back"
	ExecuteProcess       = "execute"
	PrintCommands        = "print_commands"
)

type UserInteraction interface {
	NextInteraction(string) UserInteraction
	Execute() bool
	setNextInteraction(command string, nextProcess UserInteraction)
	GetName() string
	getNextInteraction() UserInteraction
	PrintPossibleInteraction(int)
	PrintCommands(int)
	getReleaseProcess() ProcessInterface
}

type userInteraction struct {
	Name            string
	interactions    []UserInteraction
	nextInteraction UserInteraction
	releaseProcess  ProcessInterface
}

var (
	MainInteraction UserInteraction
)

func init() {
	MainInteraction = &userInteraction{Name: NICMProcess}
	initProcesses()
}

func initProcesses() {
	initClientProcess()
	initServerProcess()
	initNIGProcess()
	initGearsProcess()
	initDTSTestProcess()
	initDTSProdProcess()
}

func initDefaultProcesses(process, next UserInteraction) {
	executeClientProcess := userInteraction{Name: ExecuteProcess}
	process.setNextInteraction(ExecuteProcess, &executeClientProcess)

	process.setNextInteraction(PrintCommands, next)
	process.setNextInteraction(PreviousProcess, next)
	process.setNextInteraction(HelpCommandText, next)
	process.setNextInteraction(ExitCommandText, next)
}

func initClientProcess() {
	clientReleaseProcess := userInteraction{Name: ClientReleaseProcess, releaseProcess: NewClient()}

	MainInteraction.setNextInteraction(ClientReleaseProcess, &clientReleaseProcess)

	initDefaultProcesses(&clientReleaseProcess, MainInteraction)
}

func initServerProcess() {
	serverReleaseProcess := userInteraction{Name: ServerReleaseProcess, releaseProcess: NewServer()}

	MainInteraction.setNextInteraction(ServerReleaseProcess, &serverReleaseProcess)

	initDefaultProcesses(&serverReleaseProcess, MainInteraction)
}

func initNIGProcess() {
	nigReleaseProcess := userInteraction{Name: NIGReleaseProcess, releaseProcess: NewNIG()}

	MainInteraction.setNextInteraction(NIGReleaseProcess, &nigReleaseProcess)

	initDefaultProcesses(&nigReleaseProcess, MainInteraction)
}

func initGearsProcess() {
	gearsReleaseProcess := userInteraction{Name: GearsReleaseProcess, releaseProcess: NewGears()}

	MainInteraction.setNextInteraction(GearsReleaseProcess, &gearsReleaseProcess)

	initDefaultProcesses(&gearsReleaseProcess, MainInteraction)
}

func initDTSTestProcess() {
	dtsProcess := userInteraction{Name: DTSTestRelease, releaseProcess: NewDTSTestServer()}

	MainInteraction.setNextInteraction(DTSTestRelease, &dtsProcess)

	initDefaultProcesses(&dtsProcess, MainInteraction)
}

func initDTSProdProcess() {
	dtsProcess := userInteraction{Name: DTSProdRelease, releaseProcess: NewDTSProdServer()}

	MainInteraction.setNextInteraction(DTSProdRelease, &dtsProcess)

	initDefaultProcesses(&dtsProcess, MainInteraction)
}

func (up *userInteraction) setNextInteraction(command string, nextProcess UserInteraction) {
	userCommand := userInteraction{Name: command, nextInteraction: nextProcess}
	up.interactions = append(up.interactions, &userCommand)
}

func (up *userInteraction) PrintPossibleInteraction(tabs int) {
	for _, pr := range up.interactions {
		fmt.Print(strings.Repeat(" ", tabs))
		color.Yellow(pr.GetName())
	}
	fmt.Printf("\n")
}

func (up *userInteraction) NextInteraction(cmd string) UserInteraction {
	for _, pr := range up.interactions {
		if pr.GetName() == cmd {
			return pr.getNextInteraction()
		}
	}
	color.Red("next interaction not found...\n")
	return up
}

func (up *userInteraction) Execute() bool {
	if up.releaseProcess == nil {
		return false
	}
	if err := up.releaseProcess.Execute(); err != nil {
		color.Red("error executing interactions %s: \n", err.Error())
		panic(err)
	}
	return true
}

func (up *userInteraction) GetName() string {
	return up.Name
}

func (up *userInteraction) getNextInteraction() UserInteraction {
	return up.nextInteraction
}

func (up *userInteraction) PrintCommands(tabs int) {
	if up.releaseProcess != nil {
		up.releaseProcess.PrintCommands(tabs)
	}
}

func (up *userInteraction) getReleaseProcess() ProcessInterface {
	return up.releaseProcess
}

func (up *userInteraction) getRealReleaseProcess() ProcessInterface {
	prevProcess := up.NextInteraction(PreviousProcess)
	max := 20
	count := 0
	for {
		if count == max {
			return nil
		}
		count++
		if prevProcess.getReleaseProcess() != nil {
			return prevProcess.getReleaseProcess()
		} else {
			prevProcess = prevProcess.NextInteraction(PreviousProcess)
		}
	}
}
