package gears

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

var (
	ReleaseProcess = &gearsRelease{}
)

type gearsRelease struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

func (g *gearsRelease) Execute() error {
	for _, cmd := range g.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *gearsRelease) Init() {
	g.initOptions()
	g.initCommands()
}

func (g *gearsRelease) initOptions() {
	cfgMap := utils.GetConfig()
	g.Options = options.New()

	clientMap, ok := cfgMap["gears"]
	if ok {
		utils.SetOptionPaths(g.Options, clientMap)
	}
}

func (g *gearsRelease) initCommands() {
	g.commands = []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, g.Options.GetGitPath()),
		commands.New("build jars", utils.BuildJars, g.Options.GetBuildPath()),
	}
}

func (g *gearsRelease) PrintCommands(tabs int) {
	for _, c := range g.commands {
		c.Print(tabs)
	}
}

func (g *gearsRelease) PrintOptions(tabs int) {
	g.Options.Print(tabs)
}

func (g *gearsRelease) SetWorkingPath(path string) {
	g.Options.SetWorkingPath(path)
}

func (g *gearsRelease) SetGitPath(path string) {
	g.Options.SetGitPath(path)
}
