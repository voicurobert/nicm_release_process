package local

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/utils"
)

type localReleaseProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

const (
	workingPath          = "C:\\sw\\nicm\\"
	archiveName          = "nicm_products_client.zip"
	fileExtension        = "magikc"
	taskStatusScriptPath = "C:\\sw\\scripts\\run_ps.bat"
	disableTask          = "false"
	enableTask           = "true"
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
	l.Options = options.New(workingPath)
}

func getClientDirsToArchive() []string {
	return []string{"nicm_master\\nicm_products\\"}
}

func getServerDirsToArchive() []string {
	return []string{"nicm_master\\nicm_products\\", "nicm_master\\dynamic_patches", "externals\\diagnostics_mysql_151"}
}

func getDirsToSkipArchive() []string {
	return []string{"nicm_night_scripts", "nicm_nig"}
}

func getImageNames() []string {
	return []string{"nicm_open", "nicm_closed"}
}

func (l *localReleaseProcess) initCommands() {
	l.commands = []commands.CommandInterface{
		commands.NewCommand("disable Scheduled Task", utils.SetTaskStatus, taskStatusScriptPath, disableTask),
		commands.NewCommand("execute git pull", utils.ExecuteGitPull, l.Options.GetGitPath()),
		commands.NewCommand("delete magikc files", utils.DeleteFiles, l.Options.GetGitPath(), fileExtension),
		commands.NewCommand("build images", utils.BuildImages, l.Options.GetBuildPath()),
		commands.NewCommand("set writable access", utils.SetWritableAccess, l.Options.GetImagesPath(), getImageNames()),
		commands.NewCommand("creating client archive", utils.CreateArchive, l.Options.WorkingPath, archiveName, getClientDirsToArchive()),
		commands.NewCommand("creating server archive", utils.CreateArchive, l.Options.WorkingPath, archiveName, getServerDirsToArchive(), getDirsToSkipArchive()),
		commands.NewCommand("disable Scheduled Task", utils.SetTaskStatus, taskStatusScriptPath, enableTask),
	}
}

func (l *localReleaseProcess) PrintCommands() {
	for _, c := range l.commands {
		c.Print()
	}
}

func (l *localReleaseProcess) PrintOptions() {
	l.Options.Print()
}

func (l *localReleaseProcess) SetWorkingPath(path string) {
	l.Options.SetWorkingPath(path)
}

func (l *localReleaseProcess) SetBuildPath(path string) {
	l.Options.SetBuildPath(path)
}

func (l *localReleaseProcess) SetAntCommand(path string) {
	l.Options.SetAntCommand(path)
}

func (l *localReleaseProcess) SetImagesPath(path string) {
	l.Options.SetImagesPath(path)
}

func (l *localReleaseProcess) SetGitPath(path string) {
	l.Options.SetGitPath(path)
}
