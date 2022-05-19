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
	archiveName        = "nicm.zip"
	serverArchiveName  = "nicm_products.zip"
	serverArchivePath  = "/home/laur/nicm/"
	dtsArchiveName     = "nicm_products.zip"
	dtsTestArchivePath = "/dts/sw529/nicm/"
	dtsProdArchivePath = "/dts/sw529/nicm/"
)

func serverDirsToArchive() []string {
	return []string{"nicm\\nicm_products\\", "nicm\\dynamic_patches\\", "externals\\diagnostics_mysql_151\\"}
}

func NewServer() ProcessInterface {
	opts := options.NewOptionsWithConfigName("server_release")
	list := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, opts.WorkingPath, serverArchiveName, serverDirsToArchive()),
		commands.New("move archive to server", utils.MoveArchive, *opts.Server, opts.WorkingPath+serverArchiveName, serverArchivePath+serverArchiveName),
		commands.New("delete old archive", utils.DeleteThing, *opts.Server, serverArchivePath+"nicm_products_old.zip"),
		commands.New("rename archive", utils.RenameThing, *opts.Server, serverArchivePath+serverArchiveName, serverArchivePath+"nicm_products_old.zip"),
		commands.New("unzip archive", utils.Unzip, *opts.Server, serverArchivePath+serverArchiveName),
		commands.New("start job server (job.server)", utils.StartJobServerForScreenName, *opts.Server, "job.server"),
		commands.New("start job server (lni.server)", utils.StartJobServerForScreenName, *opts.Server, "lni.server"),
	}

	return newProcess(opts, list)
}

func NewDTSTestServer() ProcessInterface {
	opts := options.NewOptionsWithConfigName("dts_test_config")

	cmds := commands.NewCommandsFromConfig("dts_test_commands")

	path := opts.WorkingPath + opts.GitPath
	list := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, path, dtsArchiveName, []string{"nicm_products\\"}),
		commands.New("delete nicm_products_old", utils.DeleteThing, *opts.Server, dtsTestArchivePath+"nicm_products_old"),
		commands.New("rename nicm_products", utils.RenameThing, *opts.Server, dtsTestArchivePath+"nicm_products", dtsTestArchivePath+"nicm_products_old"),
		commands.New("delete old archive", utils.DeleteThing, *opts.Server, dtsTestArchivePath+dtsArchiveName),
		commands.New("move archive to server", utils.MoveArchive, *opts.Server, path+dtsArchiveName, dtsTestArchivePath+dtsArchiveName),
		commands.New("unzip archive", utils.Unzip, *opts.Server, dtsTestArchivePath, dtsArchiveName),
		commands.New("run commands", utils.RunDTSCommands, *opts.Server, cmds),
	}
	return newProcess(opts, list)
}

func NewDTSProdServer() ProcessInterface {
	opts := options.NewOptionsWithConfigName("dts_prod_config")
	list := []commands.CommandInterface{
		commands.New("test commands", utils.TestDTSCommands, *opts.Server),
	}

	return newProcess(opts, list)
}

func NewClient() ProcessInterface {
	opts := options.NewOptionsWithConfigName("client_release")

	list := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
		commands.New("creating archive", utils.CreateArchive, opts.WorkingPath, archiveName, []string{opts.GitPath}),
	}

	return newProcess(opts, list)
}

func NewNIG() ProcessInterface {
	opts := options.NewOptionsWithConfigName("nig")

	list := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
	}

	return newProcess(opts, list)
}

func NewGears() ProcessInterface {
	opts := options.NewOptionsWithConfigName("gears")
	list := []commands.CommandInterface{
		commands.New("git pull", utils.ExecuteGitPull, opts.GetGitPath()),
		commands.New("build jars", utils.BuildJars, opts.GetBuildPath()),
	}

	return newProcess(opts, list)
}
