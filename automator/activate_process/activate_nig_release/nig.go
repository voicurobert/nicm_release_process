package activate_nig_release

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/options"
)

type nigRelease struct {
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
	NigRelease = &nigRelease{}
)
