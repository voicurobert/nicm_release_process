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
	workingPath = "C:\\sw\\nicm\\"
	gitPath     = "nicm_master\\"
	buildPath   = "run\\nicm430\\"
	antCommand  = "build"
	imagesPath  = "images\\main"
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
	c.Options.SetWorkingPath(workingPath)
	c.Options.SetGitPath(gitPath)
	c.Options.SetBuildPath(buildPath)
	c.Options.SetAntCommand(antCommand)
	c.Options.SetImagesPath(imagesPath)
}

func (c *clientReleaseProcess) initCommands() {
	//c.commands = append(c.commands, commands.NewCommand("delete magikc files", utils.DeleteFiles, c.Options.GetGitPath()))
	//c.commands = append(c.commands, commands.NewCommand("execute git pull", utils.ExecuteGitPull, c.Options.GetGitPath()))
	//c.commands = append(c.commands, commands.NewCommand("build images", utils.BuildImages, c.Options.GetBuildPath()))
	//c.commands = append(c.commands, commands.NewCommand("set writable access", utils.SetWritableAccess, c.Options.GetImagesPath(), "nicm_open", "nicm_closed"))
	//c.commands = append(c.commands, commands.NewCommand(
	//	"creating archive",
	//	utils.CreateArchive,
	//	c.Options.WorkingPath, "nicm_products_client.zip", []string{"nicm_master\\nicm_products\\"}))
	//c.commands = append(c.commands, commands.NewCommand("disable Task", utils.ExecutePowerShell, "Disable-ScheduledTask -TaskPath \"\\NICM\\\" -TaskName \"Test\"", "-verb RunAs"))
	c.commands = append(c.commands, commands.NewCommand("disable Task", utils.SetTaskStatus, "false"))
}

func (c *clientReleaseProcess) PrintCommands() {
	for _, c := range c.commands {
		c.Print()
	}
}

func (c *clientReleaseProcess) PrintOptions() {
	c.Options.Print()
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
