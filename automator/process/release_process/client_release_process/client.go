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
	archiveName = "nicm.zip"
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
	cfgMap := utils.GetConfig()
	c.Options = options.New()

	clientMap, ok := cfgMap["client_release"]
	if ok {
		utils.SetOptionPaths(c.Options, clientMap)
	}
}

func (c *clientReleaseProcess) dirsToArchive() []string {
	return []string{c.Options.GitPath}
}

func (c *clientReleaseProcess) initCommands() {
	c.commands = nil
	c.commands = []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, c.Options.GetGitPath()),
		commands.New("build jars", utils.BuildJars, c.Options.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, c.Options.WorkingPath, archiveName, c.dirsToArchive()),
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
	c.initCommands()
}

func (c *clientReleaseProcess) SetGitPath(path string) {
	c.Options.SetGitPath(path)
	c.initCommands()
}
