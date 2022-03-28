package oro

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

type oroDtsDeployProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

const (
	//workingPath     = "C:\\NIG\\"
	workingPath    = "C:\\sw\\nicm\\nicm_529\\"
	archiveName    = "nicm.zip"
	magikExtension = ".magik"
)

var (
	ReleaseProcess = &oroDtsDeployProcess{}
)

func (d *oroDtsDeployProcess) Execute() error {
	for _, cmd := range d.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *oroDtsDeployProcess) Init() {
	d.initOptions()
	d.initCommands()
}

func (d *oroDtsDeployProcess) initOptions() {
	d.Options = options.New(workingPath)
}

func (d *oroDtsDeployProcess) dirsToArchive() []string {
	return []string{d.Options.GitPath}
}

func (d *oroDtsDeployProcess) initCommands() {
	d.commands = nil
	d.commands = []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, d.Options.GetGitPath()),
		commands.New("build jars", utils.BuildJars, d.Options.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, d.Options.WorkingPath, archiveName, d.dirsToArchive()),
	}
}

func (d *oroDtsDeployProcess) PrintCommands(tabs int) {
	for _, c := range d.commands {
		c.Print(tabs)
	}
}

func (d *oroDtsDeployProcess) PrintOptions(tabs int) {
	d.Options.Print(tabs)
}

func (d *oroDtsDeployProcess) SetWorkingPath(path string) {
	d.Options.SetWorkingPath(path)
	d.initCommands()
}

func (d *oroDtsDeployProcess) SetGitPath(path string) {
	d.Options.SetGitPath(path)
	d.initCommands()
}
