package process

import (
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/release_process/client"
	"github.com/voicurobert/nicm_release_process/automator/release_process/server"
)

const (
	HelpCommandText        = "help"
	ExitCommandText        = "exit"
	NICMProcess            = "NICM Process"
	PrepareReleaseProcess  = "prepare release"
	ClientReleaseProcess   = "client release"
	ServerReleaseProcess   = "server release"
	ActivateReleaseProcess = "activate release"
	NIGReleaseProcess      = "nig release"
	GearsReleaseProcess    = "gears release"
	PreviousProcess        = "back"
	ExecuteProcess         = "execute"
	PrintCommands          = "print commands"
	PrintOptions           = "print options"
)

var (
	Process UserProcessInterface
)

type userProcess struct {
	Name            string
	processes       []UserProcessInterface
	nextUserProcess UserProcessInterface
	internalProcess commands.ReleaseProcessInterface
}

type UserProcessInterface interface {
	NextProcess(string) UserProcessInterface
	Execute() UserProcessInterface
	setNextProcess(command string, nextProcess UserProcessInterface)
	GetName() string
	getNextProcess() UserProcessInterface
	PrintPossibleProcesses()
	PrintCommands()
	PrintOptions()
}

func init() {
	Process = &userProcess{Name: NICMProcess}
	initPrepareReleaseProcess()
	initActivateReleaseProcess()
}

func initPrepareReleaseProcess() {
	prepareReleaseProcess := userProcess{Name: PrepareReleaseProcess}

	Process.setNextProcess(PrepareReleaseProcess, &prepareReleaseProcess)

	initClientReleaseProcess(&prepareReleaseProcess)
	initServerReleaseProcess(&prepareReleaseProcess)

	prepareReleaseProcess.setNextProcess(PreviousProcess, Process)
	prepareReleaseProcess.setNextProcess(HelpCommandText, Process)
	prepareReleaseProcess.setNextProcess(ExitCommandText, Process)
}

func initClientReleaseProcess(releaseProcess UserProcessInterface) {
	clientReleaseProcess := userProcess{Name: ClientReleaseProcess, internalProcess: client.ReleaseProcess}

	releaseProcess.setNextProcess(ClientReleaseProcess, &clientReleaseProcess)

	clientReleaseProcess.internalProcess.Init()

	clientReleaseProcess.setNextProcess(PreviousProcess, releaseProcess)
	executeClientProcess := userProcess{Name: ExecuteProcess}
	clientReleaseProcess.setNextProcess(ExecuteProcess, &executeClientProcess)
	clientReleaseProcess.setNextProcess(HelpCommandText, Process)
	clientReleaseProcess.setNextProcess(ExitCommandText, Process)
}

func initServerReleaseProcess(releaseProcess UserProcessInterface) {
	serverReleaseProcess := userProcess{Name: ServerReleaseProcess, internalProcess: server.ReleaseProcess}
	releaseProcess.setNextProcess(ServerReleaseProcess, &serverReleaseProcess)
	serverReleaseProcess.setNextProcess(PreviousProcess, releaseProcess)
	executeServerProcess := userProcess{Name: ExecuteProcess}
	serverReleaseProcess.setNextProcess(ExecuteProcess, &executeServerProcess)
}

func initActivateReleaseProcess() {
	activateReleaseProcess := userProcess{Name: ActivateReleaseProcess}
	nigReleaseProcess := userProcess{Name: NIGReleaseProcess}
	gearsReleaseProcess := userProcess{Name: GearsReleaseProcess}
	activateReleaseProcess.setNextProcess(NIGReleaseProcess, &nigReleaseProcess)
	activateReleaseProcess.setNextProcess(GearsReleaseProcess, &gearsReleaseProcess)

	executeNIGRelease := userProcess{Name: ExecuteProcess}
	executeGearsRelease := userProcess{Name: ExecuteProcess}

	gearsReleaseProcess.setNextProcess(PreviousProcess, &activateReleaseProcess)
	gearsReleaseProcess.setNextProcess(ExecuteProcess, &executeGearsRelease)
	nigReleaseProcess.setNextProcess(PreviousProcess, &activateReleaseProcess)
	nigReleaseProcess.setNextProcess(ExecuteProcess, &executeNIGRelease)
	activateReleaseProcess.setNextProcess(PreviousProcess, Process)

	Process.setNextProcess(ActivateReleaseProcess, &activateReleaseProcess)
}

func (up *userProcess) setNextProcess(command string, nextProcess UserProcessInterface) {
	userCommand := userProcess{Name: command, nextUserProcess: nextProcess}
	up.processes = append(up.processes, &userCommand)
}

func (up *userProcess) PrintPossibleProcesses() {
	for _, process := range up.processes {
		fmt.Printf("\t")
		fmt.Print(process.GetName())
		fmt.Printf("\n")
	}
}

func (up *userProcess) NextProcess(cmd string) UserProcessInterface {
	for _, process := range up.processes {
		if process.GetName() == cmd {
			return process.getNextProcess()
		}
	}
	fmt.Println("Unknown command...")
	return up
}

func (up *userProcess) Execute() UserProcessInterface {
	if up.internalProcess == nil {
		return up.NextProcess(PreviousProcess)
	}
	if err := up.internalProcess.Execute(); err != nil {
		fmt.Printf("error executing process %s:", err.Error())
		panic(err)
	}
	return up.NextProcess(PreviousProcess)
}

func (up *userProcess) GetName() string {
	return up.Name
}

func (up *userProcess) getNextProcess() UserProcessInterface {
	return up.nextUserProcess
}

func (up *userProcess) PrintCommands() {
	if up.internalProcess != nil {
		up.internalProcess.PrintCommands()
	}
}

func (up *userProcess) PrintOptions() {
	if up.internalProcess != nil {
		up.internalProcess.PrintOptions()
	}
}
