package server

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/options"
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
