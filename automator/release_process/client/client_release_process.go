package client

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/options"
	"github.com/voicurobert/nicm_release_process/utils"
)

type clientReleaseProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

const (
	workingPath        = "C:\\sw\\nicm\\"
	gitPath            = workingPath + "nicm_master\\"
	deleteFilesCommand = "C:\\sw\\demo\\"
	buildPath          = "nicm_master\\run\\nicm430"
	antCommand         = "build"
	imagesPath         = "nicm_master\\run\\nicm430\\images\\main"
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
	c.Options = &options.Options{}
	c.Options.SetWorkingPath("C:\\sw\\nicm\\")
	c.Options.SetGitPath("nicm_master\\")
	c.Options.SetBuildPath("run\\nicm430")
	c.Options.SetAntCommand("build")
	c.Options.SetImagesPath("images\\main")
}

func (c *clientReleaseProcess) initCommands() {
	c.commands = append(c.commands, commands.NewCommand("delete magikc files", utils.DeleteFiles, c.Options.GetGitPath()))
	c.commands = append(c.commands, commands.NewCommand("execute git pull", utils.ExecuteGitPull, c.Options.GetGitPath()))
	c.commands = append(c.commands, commands.NewCommand("build images", utils.BuildImages, c.Options.GetBuildPath(), c.Options.AntCommand))
	c.commands = append(c.commands, commands.NewCommand("set writable access", utils.SetWritableAccess, c.Options.GetImagesPath(), "nicm_open", "nicm_closed"))
	c.commands = append(c.commands, commands.NewCommand("creating archive", utils.CreateArchive, c.Options.WorkingPath, "nicm_products.zip", "nicm_master\\nicm_products\\", "externals\\diagnostics_mysql_151"))
}

func (c *clientReleaseProcess) PrintCommands() {
	for _, c := range c.commands {
		c.Print()
	}
}

func (c *clientReleaseProcess) PrintOptions() {
	c.Options.Print()
}
