package server

import (
	"github.com/voicurobert/nicm_release_process/automator/process/commands"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

const (
	workingPath       = "C:\\sw\\nicm\\"
	archiveName       = "nicm_products.zip"
	serverArchivePath = "/home/laur/nicm/"
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
	return []string{"nicm\\nicm_products\\", "nicm\\dynamic_patches\\", "externals\\diagnostics_mysql_151\\"}
}

func getDirsToSkipArchive() []string {
	return []string{"nicm_night_scripts", "nicm_nig"}
}

func (s *serverReleaseProcess) initCommands() {
	s.commands = []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, s.Options.GetGitPath()),
		commands.New("build jars", utils.BuildJars, s.Options.GetBuildPath()),
		commands.New("delete magik files", utils.DeleteFiles, s.Options.GetGitPath(), ".magik"),
		//commands.New("creating archive", utils.CreateArchive, s.Options.WorkingPath, archiveName, getDirsToArchive(), getDirsToSkipArchive()),
		//commands.New("move archive to server", utils.MoveArchive, s.Options.WorkingPath+archiveName, serverArchivePath+archiveName),
		commands.New("delete old archive", utils.DeleteOldArchive, serverArchivePath+"nicm_products_old.zip"),
		commands.New("rename archive", utils.RenameArchive, serverArchivePath+archiveName, serverArchivePath+"nicm_products_old.zip"),
		commands.New("unzip archive", utils.Unzip, serverArchivePath+archiveName),
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

func (s *serverReleaseProcess) SetGitPath(path string) {
	s.Options.SetGitPath(path)
}
