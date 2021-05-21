package nig

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

const (
	workingPath          = "C:\\NIG\\"
	taskStatusScriptPath = "scripts\\run_ps.bat"
	disableTask          = "false"
	enableTask           = "true"

	magikcExtension = ".magikc"
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
		commands.NewCommand("disable Scheduled Task", utils.SetTaskStatus, taskStatusScriptPath, disableTask),
		commands.NewCommand("execute git pull", utils.ExecuteGitPull, n.Options.GetGitPath()),
		commands.NewCommand("delete magikc files", utils.DeleteFiles, n.Options.GetGitPath(), magikcExtension),
		commands.NewCommand("build images", utils.BuildImages, n.Options.GetBuildPath()),
		commands.NewCommand("disable Scheduled Task", utils.SetTaskStatus, taskStatusScriptPath, enableTask),
	}
}

func (n *nigRelease) PrintCommands() {
	for _, c := range n.commands {
		c.Print()
	}
}

func (n *nigRelease) PrintOptions() {
	n.Options.Print()
}

func (n *nigRelease) SetWorkingPath(path string) {
	n.Options.SetWorkingPath(path)
}

func (n *nigRelease) SetBuildPath(path string) {
	n.Options.SetBuildPath(path)
}

func (n *nigRelease) SetAntCommand(path string) {
	n.Options.SetAntCommand(path)
}

func (n *nigRelease) SetImagesPath(path string) {
	n.Options.SetImagesPath(path)
}

func (n *nigRelease) SetGitPath(path string) {
	n.Options.SetGitPath(path)
}
