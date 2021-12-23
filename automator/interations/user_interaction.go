package interations

import (
	"fmt"
	"github.com/voicurobert/nicm_release_process/automator/process/activate_process/gears"
	"github.com/voicurobert/nicm_release_process/automator/process/activate_process/nig_release"
	local "github.com/voicurobert/nicm_release_process/automator/process/local_release"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/process/release_process/client_release_process"
	"github.com/voicurobert/nicm_release_process/automator/process/release_process/server_release_process"
	"strings"
)

const (
	HelpCommandText        = "help"
	ExitCommandText        = "exit"
	NICMProcess            = "nicm_process"
	PrepareReleaseProcess  = "prepare_release"
	ClientReleaseProcess   = "client_release"
	LocalReleaseProcess    = "local_release"
	ServerReleaseProcess   = "server_release"
	ActivateReleaseProcess = "activate_release"
	NIGReleaseProcess      = "nig_release"
	GearsReleaseProcess    = "gears_release"
	PreviousProcess        = "back"
	ExecuteProcess         = "execute"
	PrintCommands          = "print_commands"
)

var (
	MainInteraction UserInteractionInterface
)

type ReleaseProcessInterface interface {
	Execute() error
	Init()
	PrintCommands(int)
	PrintOptions(int)
	options.SetOptionsInterface
}

type userInteraction struct {
	Name            string
	interactions    []UserInteractionInterface
	nextInteraction UserInteractionInterface
	releaseProcess  ReleaseProcessInterface
}

type UserInteractionInterface interface {
	NextInteraction(string) UserInteractionInterface
	Execute() bool
	setNextInteraction(command string, nextProcess UserInteractionInterface)
	GetName() string
	getNextInteraction() UserInteractionInterface
	PrintPossibleInteraction(int)
	PrintCommands(int)
	PrintOptions(int)
	getReleaseProcess() ReleaseProcessInterface
	options.SetOptionsInterface
}

func init() {
	MainInteraction = &userInteraction{Name: NICMProcess}
	initPrepareReleaseProcess()
	initActivateReleaseProcess()
	initLocalReleaseProcess()
}

func initLocalReleaseProcess() {
	localReleaseProcess := userInteraction{Name: LocalReleaseProcess, releaseProcess: local.ReleaseProcess}
	MainInteraction.setNextInteraction(LocalReleaseProcess, &localReleaseProcess)
	localReleaseProcess.releaseProcess.Init()
	executeClientProcess := userInteraction{Name: ExecuteProcess}
	localReleaseProcess.setNextInteraction(ExecuteProcess, &executeClientProcess)
	localReleaseProcess.setNextInteraction(PrintCommands, &localReleaseProcess)
	initOptionsProcess(&localReleaseProcess)
	initDefaultProcesses(&localReleaseProcess, MainInteraction)
}

func initPrepareReleaseProcess() {
	prepareReleaseProcess := userInteraction{Name: PrepareReleaseProcess}

	MainInteraction.setNextInteraction(PrepareReleaseProcess, &prepareReleaseProcess)

	initClientReleaseProcess(&prepareReleaseProcess)
	initServerReleaseProcess(&prepareReleaseProcess)

	initDefaultProcesses(&prepareReleaseProcess, MainInteraction)
}

func initDefaultProcesses(process, next UserInteractionInterface) {
	process.setNextInteraction(PreviousProcess, next)
	process.setNextInteraction(HelpCommandText, next)
	process.setNextInteraction(ExitCommandText, next)
}

func initClientReleaseProcess(releaseProcess UserInteractionInterface) {
	clientReleaseProcess := userInteraction{Name: ClientReleaseProcess, releaseProcess: client.ReleaseProcess}

	releaseProcess.setNextInteraction(ClientReleaseProcess, &clientReleaseProcess)

	clientReleaseProcess.releaseProcess.Init()

	executeClientProcess := userInteraction{Name: ExecuteProcess}
	clientReleaseProcess.setNextInteraction(ExecuteProcess, &executeClientProcess)
	clientReleaseProcess.setNextInteraction(PrintCommands, releaseProcess)
	initOptionsProcess(&clientReleaseProcess)
	initDefaultProcesses(&clientReleaseProcess, releaseProcess)
}

func initOptionsProcess(process UserInteractionInterface) {
	process.setNextInteraction(options.PrintOptions, MainInteraction)
	setOptionsProcess := userInteraction{Name: options.SetOptions}
	process.setNextInteraction(options.SetOptions, &setOptionsProcess)
	initSetOptionsProcess(&setOptionsProcess)
	initDefaultProcesses(&setOptionsProcess, process)
}

func initSetOptionsProcess(optionsProcess UserInteractionInterface) {
	setWorkingPathProcess := userInteraction{Name: options.SetWorkingPath}
	optionsProcess.setNextInteraction(options.SetWorkingPath, &setWorkingPathProcess)
	initDefaultProcesses(&setWorkingPathProcess, optionsProcess)

	setGitPathProcess := userInteraction{Name: options.SetGitPath}
	optionsProcess.setNextInteraction(options.SetGitPath, &setGitPathProcess)
	initDefaultProcesses(&setGitPathProcess, optionsProcess)
}

func initServerReleaseProcess(releaseProcess UserInteractionInterface) {
	serverReleaseProcess := userInteraction{Name: ServerReleaseProcess, releaseProcess: server.ReleaseProcess}

	releaseProcess.setNextInteraction(ServerReleaseProcess, &serverReleaseProcess)

	serverReleaseProcess.releaseProcess.Init()

	executeServerProcess := userInteraction{Name: ExecuteProcess}
	serverReleaseProcess.setNextInteraction(ExecuteProcess, &executeServerProcess)
	initOptionsProcess(&serverReleaseProcess)
	serverReleaseProcess.setNextInteraction(PrintCommands, releaseProcess)
	initDefaultProcesses(&serverReleaseProcess, releaseProcess)
}

func initActivateReleaseProcess() {
	activateReleaseProcess := userInteraction{Name: ActivateReleaseProcess}

	initNIGReleaseProcess(&activateReleaseProcess)
	initGearsReleaseProcess(&activateReleaseProcess)

	initDefaultProcesses(&activateReleaseProcess, MainInteraction)
	MainInteraction.setNextInteraction(ActivateReleaseProcess, &activateReleaseProcess)
}

func initNIGReleaseProcess(process UserInteractionInterface) {
	nigReleaseProcess := userInteraction{Name: NIGReleaseProcess, releaseProcess: nig.ReleaseProcess}
	nigReleaseProcess.releaseProcess.Init()
	process.setNextInteraction(NIGReleaseProcess, &nigReleaseProcess)
	executeNIGRelease := userInteraction{Name: ExecuteProcess}
	nigReleaseProcess.setNextInteraction(ExecuteProcess, &executeNIGRelease)
	nigReleaseProcess.setNextInteraction(PrintCommands, process)
	initOptionsProcess(&nigReleaseProcess)
	initDefaultProcesses(&nigReleaseProcess, process)
}

func initGearsReleaseProcess(process UserInteractionInterface) {
	gearsReleaseProcess := userInteraction{Name: GearsReleaseProcess, releaseProcess: gears.ReleaseProcess}
	gearsReleaseProcess.releaseProcess.Init()
	process.setNextInteraction(GearsReleaseProcess, &gearsReleaseProcess)

	executeGearsRelease := userInteraction{Name: ExecuteProcess}
	gearsReleaseProcess.setNextInteraction(ExecuteProcess, &executeGearsRelease)
	gearsReleaseProcess.setNextInteraction(PrintCommands, process)
	initOptionsProcess(&gearsReleaseProcess)
	initDefaultProcesses(&gearsReleaseProcess, process)
}

func (up *userInteraction) setNextInteraction(command string, nextProcess UserInteractionInterface) {
	userCommand := userInteraction{Name: command, nextInteraction: nextProcess}
	up.interactions = append(up.interactions, &userCommand)
}

func (up *userInteraction) PrintPossibleInteraction(tabs int) {
	for _, process := range up.interactions {
		fmt.Printf(strings.Repeat(" ", tabs))
		fmt.Print(process.GetName())
		fmt.Printf("\n")
	}
}

func (up *userInteraction) NextInteraction(cmd string) UserInteractionInterface {
	for _, process := range up.interactions {
		if process.GetName() == cmd {
			return process.getNextInteraction()
		}
	}
	fmt.Println("next interaction not found...")
	return up
}

func (up *userInteraction) Execute() bool {
	if up.releaseProcess == nil {
		return false
	}
	if err := up.releaseProcess.Execute(); err != nil {
		fmt.Printf("error executing interations %s:", err.Error())
		panic(err)
	}
	return true
}

func (up *userInteraction) GetName() string {
	return up.Name
}

func (up *userInteraction) getNextInteraction() UserInteractionInterface {
	return up.nextInteraction
}

func (up *userInteraction) PrintCommands(tabs int) {
	if up.releaseProcess != nil {
		up.releaseProcess.PrintCommands(tabs)
	}
}

func (up *userInteraction) PrintOptions(tabs int) {
	if up.releaseProcess != nil {
		up.releaseProcess.PrintOptions(tabs)
	}
}

func (up *userInteraction) SetWorkingPath(path string) {
	p := up.getRealReleaseProcess()
	p.SetWorkingPath(path)
}

func (up *userInteraction) SetGitPath(path string) {
	p := up.getRealReleaseProcess()
	p.SetGitPath(path)
}

func (up *userInteraction) getReleaseProcess() ReleaseProcessInterface {
	return up.releaseProcess
}

func (up *userInteraction) getRealReleaseProcess() ReleaseProcessInterface {
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
