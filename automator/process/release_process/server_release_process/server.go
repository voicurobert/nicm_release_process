package server

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

const (
	workingPath       = "C:\\sw\\nicm\\"
	archiveName       = "nicm_products_server.zip"
	magikcExtension   = ".magikc"
	serverArchivePath = "/nicm/"
)

type serverReleaseProcess struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

var (
	ReleaseProcess = &serverReleaseProcess{}
)

func (s *serverReleaseProcess) Execute() error {
	for _, cmd := range s.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serverReleaseProcess) Init() {
	s.initOptions()
	s.initCommands()
}

func (s *serverReleaseProcess) initOptions() {
	s.Options = options.New(workingPath)
}

func getDirsToArchive() []string {
	return []string{"nicm_master\\nicm_products\\", "nicm_master\\dynamic_patches\\", "externals\\diagnostics_mysql_151\\"}
}

func getDirsToSkipArchive() []string {
	return []string{"nicm_night_scripts", "nicm_nig"}
}

func getImageNames() []string {
	return []string{"nicm_open", "nicm_closed"}
}

func (s *serverReleaseProcess) initCommands() {
	s.commands = []commands.CommandInterface{
		//commands.NewCommand("delete magikc files", utils.DeleteFiles, s.Options.GetGitPath(), magikcExtension),
		//commands.NewCommand("execute git pull", utils.ExecuteGitPull, s.Options.GetGitPath()),
		//commands.NewCommand("build images", utils.BuildImages, s.Options.GetBuildPath(), s.Options.AntCommand),
		//commands.NewCommand("set writable access", utils.SetWritableAccess, s.Options.GetImagesPath(), getImageNames()),
		//commands.NewCommand("creating archive", utils.CreateArchive, s.Options.WorkingPath, archiveName, getDirsToArchive(), getDirsToSkipArchive()),
		commands.NewCommand("move archive to server", utils.MoveArchive, s.Options.WorkingPath, serverArchivePath, archiveName),
	}
}

func (s *serverReleaseProcess) PrintCommands(tabs int) {
	for _, s := range s.commands {
		s.Print(tabs)
	}
}

func (s *serverReleaseProcess) PrintOptions(tabs int) {

}

func (s *serverReleaseProcess) SetWorkingPath(path string) {
	s.Options.SetWorkingPath(path)
}

func (s *serverReleaseProcess) SetBuildPath(path string) {
	s.Options.SetBuildPath(path)
}

func (s *serverReleaseProcess) SetAntCommand(path string) {
	s.Options.SetAntCommand(path)
}

func (s *serverReleaseProcess) SetImagesPath(path string) {
	s.Options.SetImagesPath(path)
}

func (s *serverReleaseProcess) SetGitPath(path string) {
	s.Options.SetGitPath(path)
}
