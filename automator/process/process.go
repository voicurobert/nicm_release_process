package process

import (
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/options"
	"github.com/voicurobert/nicm_release_process/automator/release_process/client"
	"github.com/voicurobert/nicm_release_process/automator/release_process/server"
)

const (
	HelpCommandText        = "help"
	ExitCommandText        = "exit"
	NICMProcess            = "nicm_process"
	PrepareReleaseProcess  = "prepare_release"
	ClientReleaseProcess   = "client_release"
	ServerReleaseProcess   = "server_release"
	ActivateReleaseProcess = "activate_release"
	NIGReleaseProcess      = "nig_release"
	GearsReleaseProcess    = "gears_release"
	PreviousProcess        = "back"
	ExecuteProcess         = "execute"
	PrintCommands          = "print_commands"
)

var (
	Process UserProcessInterface
)

type ReleaseProcessInterface interface {
	Execute() error
	Init()
	PrintCommands()
	PrintOptions()
	options.SetOptionsInterface
}

type userProcess struct {
	Name            string
	processes       []UserProcessInterface
	nextUserProcess UserProcessInterface
	internalProcess ReleaseProcessInterface
}

type UserProcessInterface interface {
	NextProcess(string) UserProcessInterface
	Execute() bool
	setNextProcess(command string, nextProcess UserProcessInterface)
	GetName() string
	getNextProcess() UserProcessInterface
	PrintPossibleProcesses()
	PrintCommands()
	PrintOptions()
	getInternalProcess() ReleaseProcessInterface
	options.SetOptionsInterface
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

	initDefaultProcesses(&prepareReleaseProcess, Process)
}

func initDefaultProcesses(process, next UserProcessInterface) {
	process.setNextProcess(PreviousProcess, next)
	process.setNextProcess(HelpCommandText, next)
	process.setNextProcess(ExitCommandText, next)
}

func initClientReleaseProcess(releaseProcess UserProcessInterface) {
	clientReleaseProcess := userProcess{Name: ClientReleaseProcess, internalProcess: client.ReleaseProcess}

	releaseProcess.setNextProcess(ClientReleaseProcess, &clientReleaseProcess)

	clientReleaseProcess.internalProcess.Init()

	executeClientProcess := userProcess{Name: ExecuteProcess}
	clientReleaseProcess.setNextProcess(ExecuteProcess, &executeClientProcess)
	initOptionsProcess(&clientReleaseProcess)
	initDefaultProcesses(&clientReleaseProcess, releaseProcess)
}

func initOptionsProcess(process UserProcessInterface) {
	process.setNextProcess(options.PrintOptions, Process)
	setOptionsProcess := userProcess{Name: options.SetOptions}
	process.setNextProcess(options.SetOptions, &setOptionsProcess)
	initSetOptionsProcess(&setOptionsProcess)
	initDefaultProcesses(&setOptionsProcess, process)
}

func initSetOptionsProcess(optionsProcess UserProcessInterface) {
	setWorkingPathProcess := userProcess{Name: options.SetWorkingPath}
	optionsProcess.setNextProcess(options.SetWorkingPath, &setWorkingPathProcess)
	initDefaultProcesses(&setWorkingPathProcess, optionsProcess)

	setGitPathProcess := userProcess{Name: options.SetGitPath}
	optionsProcess.setNextProcess(options.SetGitPath, &setGitPathProcess)
	initDefaultProcesses(&setGitPathProcess, optionsProcess)

	setBuildPathProcess := userProcess{Name: options.SetBuildPath}
	optionsProcess.setNextProcess(options.SetBuildPath, &setBuildPathProcess)
	initDefaultProcesses(&setBuildPathProcess, optionsProcess)

	setAntCommandProcess := userProcess{Name: options.SetAntCommand}
	optionsProcess.setNextProcess(options.SetAntCommand, &setAntCommandProcess)
	initDefaultProcesses(&setAntCommandProcess, optionsProcess)

	setImagesCommandProcess := userProcess{Name: options.SetImagesPath}
	optionsProcess.setNextProcess(options.SetImagesPath, &setImagesCommandProcess)
	initDefaultProcesses(&setImagesCommandProcess, optionsProcess)
}

func initServerReleaseProcess(releaseProcess UserProcessInterface) {
	serverReleaseProcess := userProcess{Name: ServerReleaseProcess, internalProcess: server.ReleaseProcess}

	releaseProcess.setNextProcess(ServerReleaseProcess, &serverReleaseProcess)

	serverReleaseProcess.internalProcess.Init()

	executeServerProcess := userProcess{Name: ExecuteProcess}
	serverReleaseProcess.setNextProcess(ExecuteProcess, &executeServerProcess)
	serverReleaseProcess.setNextProcess(options.PrintOptions, releaseProcess)
	initDefaultProcesses(&serverReleaseProcess, releaseProcess)
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

func (up *userProcess) Execute() bool {
	if up.internalProcess == nil {
		return false
	}
	if err := up.internalProcess.Execute(); err != nil {
		fmt.Printf("error executing process %s:", err.Error())
		panic(err)
	}
	return true
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

func (up *userProcess) getInternalProcess() ReleaseProcessInterface {
	return up.internalProcess
}

func (up *userProcess) SetWorkingPath(path string) {
	p := up.getReleaseProcessByName(ClientReleaseProcess)
	p.SetWorkingPath(path)
}

func (up *userProcess) SetGitPath(path string) {
	p := up.getReleaseProcessByName(ClientReleaseProcess)
	p.SetGitPath(path)
}

func (up *userProcess) SetBuildPath(path string) {
	p := up.getReleaseProcessByName(ClientReleaseProcess)
	p.SetBuildPath(path)
}

func (up *userProcess) SetImagesPath(path string) {
	p := up.getReleaseProcessByName(ClientReleaseProcess)
	p.SetImagesPath(path)
}

func (up *userProcess) SetAntCommand(path string) {
	p := up.getReleaseProcessByName(ClientReleaseProcess)
	p.SetAntCommand(path)
}

func (up *userProcess) getReleaseProcessByName(name string) ReleaseProcessInterface {
	prevProcess := up.NextProcess(PreviousProcess)
	for {
		if prevProcess.GetName() == name {
			return prevProcess.getInternalProcess()
		} else {
			prevProcess = prevProcess.NextProcess(PreviousProcess)
		}
	}
}
