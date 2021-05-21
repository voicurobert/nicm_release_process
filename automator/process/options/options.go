package options

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	PrintOptions   = "print_options"
	SetOptions     = "set_options"
	SetWorkingPath = "set_working_path"
	SetGitPath     = "set_git_path"
	SetBuildPath   = "set_build_path"
	SetAntCommand  = "set_ant_command"
	SetImagesPath  = "set_images_path"

	SetOptionSeparator = "="

	gitPath      = "nicm_master\\"
	buildPath    = "run\\nicm430\\"
	antCommand   = "build"
	imagesPath   = "images\\"
	buildEnvPath = "nicm_products\\nicm_build\\"
)

type SetOptionsInterface interface {
	SetWorkingPath(path string)
	SetGitPath(path string)
	SetBuildPath(path string)
	SetAntCommand(path string)
	SetImagesPath(path string)
}

type Options struct {
	WorkingPath string
	GitPath     string
	BuildPath   string
	AntCommand  string
	ImagesPath  string
}

func New(workingPath string) *Options {
	return &Options{
		WorkingPath: workingPath,
		GitPath:     gitPath,
		BuildPath:   buildPath,
		AntCommand:  antCommand,
		ImagesPath:  imagesPath,
	}
}

func (o *Options) SetWorkingPath(path string) {
	o.WorkingPath = path
}

func (o *Options) SetGitPath(path string) {
	o.GitPath = path
}

func (o *Options) SetBuildPath(path string) {
	o.BuildPath = path
}

func (o *Options) SetAntCommand(path string) {
	o.AntCommand = path
}

func (o *Options) SetImagesPath(path string) {
	o.ImagesPath = path
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

func (o *Options) GetImagesPath() string {
	var sb strings.Builder
	sb.WriteString(o.GetBuildPath())
	sb.WriteString(o.ImagesPath)
	return sb.String()
}

func (o *Options) GetBuildEnvPath() string {
	var sb strings.Builder
	sb.WriteString(o.GetGitPath())
	sb.WriteString(buildEnvPath)
	return sb.String()
}

func (o *Options) Print() {
	attrs := reflect.ValueOf(o).Elem()
	for i := 0; i < attrs.NumField(); i++ {
		name := attrs.Type().Field(i).Name
		value := attrs.Field(i).Interface()
		fmt.Printf("\t%s -> %s \n", name, value)
	}
}
