package gears

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

const (
	workingPath = "C:\\nicm_gears\\"
	gitPath     = "nicm_modules\\"
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
	g.Options = options.New(workingPath)
	g.Options.SetGitPath(gitPath)
}

func (g *gearsRelease) initCommands() {
	g.commands = []commands.CommandInterface{
		commands.NewCommand("execute git pull", utils.ExecuteGitPull, g.Options.GetGitPath()),
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
