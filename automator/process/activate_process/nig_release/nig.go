package nig

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
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
	cfgMap := utils.GetConfig()
	n.Options = options.New()

	clientMap, ok := cfgMap["nig"]
	if ok {
		utils.SetOptionPaths(n.Options, clientMap)
	}
}

func (n *nigRelease) initCommands() {
	n.commands = []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, n.Options.GetGitPath()),
		commands.New("build jars", utils.BuildJars, n.Options.GetBuildPath()),
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
