package options

import (
	"github.com/voicurobert/nicm_release_process/automator/config"
	"strings"
)

const (
	buildPath = "run5\\nicm\\"
)

type Options struct {
	WorkingPath string
	GitPath     string
	BuildPath   string
}

func NewOptionsWithConfigName(cfgName string) *Options {
	opts := New()
	cfgMap := config.GetConfig()
	clientMap, ok := cfgMap[cfgName]
	if ok {
		setPaths(opts, clientMap)
	}
	return opts
}

func setPaths(options *Options, cfgMap map[string]string) {
	if path, ok := cfgMap["git_path"]; ok {
		options.SetGitPath(strings.TrimSpace(path))
	}
	if path, ok := cfgMap["working_path"]; ok {
		options.SetWorkingPath(strings.TrimSpace(path))
	}
}

func New() *Options {
	return &Options{
		WorkingPath: "",
		GitPath:     "",
		BuildPath:   buildPath,
	}
}

func (o *Options) SetWorkingPath(path string) {
	o.WorkingPath = path
}

func (o *Options) SetGitPath(path string) {
	o.GitPath = path
}

func (o *Options) GetGitPath() string {
	var sb strings.Builder
	sb.WriteString(o.WorkingPath)
	sb.WriteString(o.GitPath)
	return sb.String()
}

func (o *Options) GetBuildPath() string {
	var sb strings.Builder
	sb.WriteString(o.GetGitPath())
	sb.WriteString(o.BuildPath)
	return sb.String()
}
