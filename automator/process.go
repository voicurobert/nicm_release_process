package automator

import (
	"github.com/voicurobert/nicm_release_process/automator/commands"
	"github.com/voicurobert/nicm_release_process/automator/options"
	"github.com/voicurobert/nicm_release_process/automator/utils"
)

type ProcessInterface interface {
	Execute() error
	PrintCommands(int)
}

type process struct {
	commands []commands.CommandInterface
	Options  *options.Options
}

func (r *process) Execute() error {
	for _, cmd := range r.commands {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func newProcess(opts *options.Options, commands []commands.CommandInterface) *process {
	return &process{
		commands: commands,
		Options:  opts,
	}
}

func (r *process) PrintCommands(tabs int) {
	for _, c := range r.commands {
		c.Print(tabs)
	}
}

const (
	archiveName       = "nicm.zip"
	serverArchiveName = "nicm_products.zip"
	serverArchivePath = "/home/laur/nicm/"
)

func serverDirsToArchive() []string {
	return []string{"nicm\\nicm_products\\", "nicm\\dynamic_patches\\", "externals\\diagnostics_mysql_151\\"}
}

func NewServer() ProcessInterface {
	opts := options.NewOptionsWithConfigName("server_release")
	cmds := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, opts.WorkingPath, serverArchiveName, serverDirsToArchive()),
		commands.New("move archive to server", utils.MoveArchive, opts.WorkingPath+serverArchiveName, serverArchivePath+serverArchiveName),
		commands.New("delete old archive", utils.DeleteOldArchive, serverArchivePath+"nicm_products_old.zip"),
		commands.New("rename archive", utils.RenameArchive, serverArchivePath+serverArchiveName, serverArchivePath+"nicm_products_old.zip"),
		commands.New("unzip archive", utils.Unzip, serverArchivePath+serverArchiveName),
		commands.New("start job server (job.server)", utils.StartJobServerForScreenName, "job.server"),
		commands.New("start job server (lni.server)", utils.StartJobServerForScreenName, "lni.server"),
	}

	return newProcess(opts, cmds)
}

func NewClient() ProcessInterface {
	opts := options.NewOptionsWithConfigName("client_release")

	cmds := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, opts.WorkingPath, archiveName, []string{opts.GitPath}),
	}

	return newProcess(opts, cmds)
}

func NewNIG() ProcessInterface {
	opts := options.NewOptionsWithConfigName("nig")

	cmds := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
	}

	return newProcess(opts, cmds)
}

func NewGears() ProcessInterface {
	opts := options.NewOptionsWithConfigName("gears")
	cmds := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
	}

	return newProcess(opts, cmds)
}
