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
	workingPath     = "C:\\NIG\\"
	archiveName     = "nicm_products_client.zip"
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

func getDirsToArchive() []string {
	return []string{"nicm_master\\nicm_products\\"}
}

func getImageNames() []string {
	return []string{"nicm_open", "nicm_closed"}
}

func (c *clientReleaseProcess) initCommands() {
	c.commands = []commands.CommandInterface{
		commands.NewCommand("execute git pull", utils.ExecuteGitPull, c.Options.GetGitPath()),
		commands.NewCommand("delete magikc files", utils.DeleteFiles, c.Options.GetGitPath(), magikcExtension),
		commands.NewCommand("build images", utils.BuildImages, c.Options.GetBuildPath()),
		commands.NewCommand("set writable access", utils.SetWritableAccess, c.Options.GetImagesPath(), getImageNames()),
		commands.NewCommand("delete magik files", utils.DeleteFiles, c.Options.GetGitPath(), magikExtension),
		commands.NewCommand("creating archive", utils.CreateArchive, c.Options.WorkingPath, archiveName, getDirsToArchive()),
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

func (c *clientReleaseProcess) SetBuildPath(path string) {
	c.Options.SetBuildPath(path)
}

func (c *clientReleaseProcess) SetAntCommand(path string) {
	c.Options.SetAntCommand(path)
}

func (c *clientReleaseProcess) SetImagesPath(path string) {
	c.Options.SetImagesPath(path)
}

func (c *clientReleaseProcess) SetGitPath(path string) {
	c.Options.SetGitPath(path)
}
