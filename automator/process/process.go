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
	Process = &userProcess{Name: NICMProcess}
)

type userProcess struct {
	Name            string
	processes       []*userProcess
	nextUserProcess *userProcess
	internalProcess commands.ReleaseProcessInterface
}

type UserProcessInterface interface {
	NextProcess(string) UserProcessInterface
	Execute() UserProcessInterface
}

func init() {
	prepareReleaseProcess := userProcess{Name: PrepareReleaseProcess}
	clientReleaseProcess := userProcess{Name: ClientReleaseProcess, internalProcess: client.ReleaseProcess}
	clientReleaseProcess.internalProcess.Init()

	serverReleaseProcess := userProcess{Name: ServerReleaseProcess, internalProcess: server.ReleaseProcess}

	prepareReleaseProcess.setNextProcess(ClientReleaseProcess, &clientReleaseProcess)
	prepareReleaseProcess.setNextProcess(ServerReleaseProcess, &serverReleaseProcess)

	clientReleaseProcess.setNextProcess(PreviousProcess, &prepareReleaseProcess)
	executeClientProcess := userProcess{Name: ExecuteProcess}
	clientReleaseProcess.setNextProcess(ExecuteProcess, &executeClientProcess)

	serverReleaseProcess.setNextProcess(PreviousProcess, &prepareReleaseProcess)
	executeServerProcess := userProcess{Name: ExecuteProcess}
	serverReleaseProcess.setNextProcess(ExecuteProcess, &executeServerProcess)

	prepareReleaseProcess.setNextProcess(PreviousProcess, Process)

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

	Process.setNextProcess(PrepareReleaseProcess, &prepareReleaseProcess)
	Process.setNextProcess(ActivateReleaseProcess, &activateReleaseProcess)
}

func (up *userProcess) setNextProcess(command string, nextProcess *userProcess) {
	userCommand := userProcess{Name: command, nextUserProcess: nextProcess}
	up.processes = append(up.processes, &userCommand)
}

func (up *userProcess) PrintPossibleProcesses() {
	for _, process := range up.processes {
		fmt.Printf("\t")
		fmt.Print(process.Name)
		fmt.Printf("\n")
	}
}

func (up *userProcess) NextProcess(cmd string) *userProcess {
	for _, process := range up.processes {
		if process.Name == cmd {
			return process.nextUserProcess
		}
	}
	fmt.Println("Unknown command...")
	return up
}

func (up *userProcess) Execute() *userProcess {
	if up.internalProcess == nil {
		return up.NextProcess(PreviousProcess)
	}
	if err := up.internalProcess.Execute(); err != nil {
		fmt.Printf("error executing process %s:", err.Error())
		panic(err)
	}
	return up.NextProcess(PreviousProcess)
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
