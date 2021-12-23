package nig

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

const (
	workingPath = "C:\\NIG\\"
	disableTask = "Disable-ScheduledTask"
	enableTask  = "Enable-ScheduledTask"

	magikcExtension             = ".magikc"
	restartGSSAgentsScriptNameA = ""
	restartGSSAgentsScriptNameB = ""
)

var (
	ReleaseProcess = &nigRelease{}
)

type nigRelease struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

func (n *nigRelease) Execute() error {
	for _, cmd := range n.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *nigRelease) Init() {
	n.initOptions()
	n.initCommands()
}

func (n *nigRelease) initOptions() {
	n.Options = options.New(workingPath)
}

func (n *nigRelease) initCommands() {
	n.commands = []commands.CommandInterface{
		commands.NewCommand("disable Scheduled Task", utils.SetScheduledTaskStatus, disableTask),
		commands.NewCommand("restart GSS agents", utils.RunPowerShellScript, n.Options.GetBuildPath()+restartGSSAgentsScriptNameA),
		commands.NewCommand("execute git pull", utils.ExecuteGitPull, n.Options.GetGitPath()),
		commands.NewCommand("delete magikc files", utils.DeleteFiles, n.Options.GetGitPath(), magikcExtension),
		commands.NewCommand("restart GSS agents", utils.RunPowerShellScript, n.Options.GetBuildPath()+restartGSSAgentsScriptNameA),
		commands.NewCommand("disable Scheduled Task", utils.SetScheduledTaskStatus, enableTask),
	}
}

func (n *nigRelease) PrintCommands(tabs int) {
	for _, c := range n.commands {
		c.Print(tabs)
	}
}

func (n *nigRelease) PrintOptions(tabs int) {
	n.Options.Print(tabs)
}

func (n *nigRelease) SetWorkingPath(path string) {
	n.Options.SetWorkingPath(path)
}

func (n *nigRelease) SetGitPath(path string) {
	n.Options.SetGitPath(path)
}
