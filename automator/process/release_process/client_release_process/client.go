package client

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

type clientReleaseProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

const (
	//workingPath     = "C:\\NIG\\"
	workingPath     = "C:\\sw\\nicm\\nicm_529\\"
	archiveName     = "nicm.zip"
	magikcExtension = ".magikc"
	magikExtension  = ".magik"
)

var (
	ReleaseProcess = &clientReleaseProcess{}
)

func (c *clientReleaseProcess) Execute() error {
	for _, cmd := range c.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *clientReleaseProcess) Init() {
	c.initOptions()
	c.initCommands()
}

func (c *clientReleaseProcess) initOptions() {
	c.Options = options.New(workingPath)
}

func (c *clientReleaseProcess) getDirsToArchive() []string {
	return []string{c.Options.GitPath}
}

func (c *clientReleaseProcess) initCommands() {
	c.commands = []commands.CommandInterface{
		//commands.NewCommand("execute git pull", utils.ExecuteGitPull, c.Options.GetGitPath()),
		//commands.NewCommand("delete magikc files", utils.DeleteFiles, c.Options.GetGitPath(), magikcExtension),
		//commands.NewCommand("delete jars", utils.DeleteJars, c.Options.GetBuildPath()),
		//commands.NewCommand("compile jars", utils.CompileJars, c.Options.GetBuildPath()),
		//commands.NewCommand("delete magik files", utils.DeleteFiles, c.Options.GetGitPath(), magikExtension),
		commands.NewCommand("creating archive", utils.CreateArchive, c.Options.WorkingPath, archiveName, c.getDirsToArchive()),
	}
}

func (c *clientReleaseProcess) PrintCommands(tabs int) {
	for _, c := range c.commands {
		c.Print(tabs)
	}
}

func (c *clientReleaseProcess) PrintOptions(tabs int) {
	c.Options.Print(tabs)
}

func (c *clientReleaseProcess) SetWorkingPath(path string) {
	c.Options.SetWorkingPath(path)
}

func (c *clientReleaseProcess) SetGitPath(path string) {
	c.Options.SetGitPath(path)
}
