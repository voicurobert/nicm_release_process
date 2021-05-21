package gears

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/utils"
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
		commands.NewCommand("build images", utils.BuildImages, g.Options.GetBuildPath()),
	}
}

func (g *gearsRelease) PrintCommands() {
	for _, c := range g.commands {
		c.Print()
	}
}

func (g *gearsRelease) PrintOptions() {
	g.Options.Print()
}

func (g *gearsRelease) SetWorkingPath(path string) {
	g.Options.SetWorkingPath(path)
}

func (g *gearsRelease) SetBuildPath(path string) {
	g.Options.SetBuildPath(path)
}

func (g *gearsRelease) SetAntCommand(path string) {
	g.Options.SetAntCommand(path)
}

func (g *gearsRelease) SetImagesPath(path string) {
	g.Options.SetImagesPath(path)
}

func (g *gearsRelease) SetGitPath(path string) {
	g.Options.SetGitPath(path)
}
