package server

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/options"
	"github.com/voicurobert/nicm_release_process/utils"
)

const (
	workingPath = "C:\\sw\\nicm\\"
	gitPath     = "nicm_master\\"
	buildPath   = "run\\nicm430"
	antCommand  = "build"
	imagesPath  = "images\\main"
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

}

func (s *serverReleaseProcess) initOptions() {
	s.Options = &options.Options{}
	s.Options.SetWorkingPath(workingPath)
	s.Options.SetGitPath(gitPath)
	s.Options.SetBuildPath(buildPath)
	s.Options.SetAntCommand(antCommand)
	s.Options.SetImagesPath(imagesPath)
}

func (s *serverReleaseProcess) initCommands() {
	s.commands = append(s.commands, commands.NewCommand("delete magikc files", utils.DeleteFiles, s.Options.GetGitPath()))
	s.commands = append(s.commands, commands.NewCommand("execute git pull", utils.ExecuteGitPull, s.Options.GetGitPath()))
	s.commands = append(s.commands, commands.NewCommand("build images", utils.BuildImages, s.Options.GetBuildPath(), s.Options.AntCommand))
	s.commands = append(s.commands, commands.NewCommand("set writable access", utils.SetWritableAccess, s.Options.GetImagesPath(), "nicm_open", "nicm_closed"))
	s.commands = append(s.commands, commands.NewCommand(
		"creating archive",
		utils.CreateArchive,
		s.Options.WorkingPath,
		"nicm_products_server.zip", []string{"nicm_master\\nicm_products\\", "nicm_master\\dynamic_patches", "externals\\diagnostics_mysql_151"}, []string{"nicm_night_scripts", "nicm_nig"}))
}

func (s *serverReleaseProcess) PrintCommands() {
	for _, s := range s.commands {
		s.Print()
	}
}

func (s *serverReleaseProcess) PrintOptions() {

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
