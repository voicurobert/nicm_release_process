package local

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

type localReleaseProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

const (
	workingPath       = "C:\\sw\\nicm\\nicm_529\\"
	clientArchiveName = "nicm_products_client.zip"
	serverArchiveName = "nicm_products_server.zip"
	magikExtension    = ".magik"
)

var (
	ReleaseProcess = &localReleaseProcess{}
)

func (l *localReleaseProcess) Execute() error {
	for _, cmd := range l.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *localReleaseProcess) Init() {
	l.initOptions()
	l.initCommands()
}

func (l *localReleaseProcess) initOptions() {
	l.Options = options.New()
}

func (l *localReleaseProcess) getClientDirsToArchive() []string {
	return []string{l.Options.GitPath}
}

func (l *localReleaseProcess) getServerDirsToArchive() []string {
	return []string{
		"nicm_master\\nicm_products\\",
		"nicm_master\\dynamic_patches",
		"externals\\diagnostics_mysql_151",
	}
}

func (l *localReleaseProcess) initCommands() {
	l.commands = []commands.CommandInterface{
		//commands.New("execute git pull", utils.ExecuteGitPull, l.Options.GetGitPath()),
		//commands.New("build jars", utils.BuildJars, l.Options.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, l.Options.WorkingPath, clientArchiveName, l.getClientDirsToArchive()),
	}
}

func (l *localReleaseProcess) PrintCommands(tabs int) {
	for _, c := range l.commands {
		c.Print(tabs)
	}
}

func (l *localReleaseProcess) PrintOptions(tabs int) {
	l.Options.Print(tabs)
}

func (l *localReleaseProcess) SetWorkingPath(path string) {
	l.Options.SetWorkingPath(path)
}

func (l *localReleaseProcess) SetGitPath(path string) {
	l.Options.SetGitPath(path)
}
