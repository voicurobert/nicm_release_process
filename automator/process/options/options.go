package options

import (
	"github.com/fatih/color"
	"reflect"
	"strings"
)

const (
	PrintOptions   = "print_options"
	SetOptions     = "set_options"
	SetWorkingPath = "set_working_path"
	SetGitPath     = "set_git_path"

	SetOptionSeparator = "="

	gitPath   = "nicm\\"
	buildPath = "run5\\nicm\\"
)

type SetOptionsInterface interface {
	SetWorkingPath(path string)
	SetGitPath(path string)
}

type Options struct {
	WorkingPath string
	GitPath     string
	BuildPath   string
}

func New(workingPath string) *Options {
	return &Options{
		WorkingPath: workingPath,
		GitPath:     gitPath,
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

func (o *Options) Print(tabs int) {
	tabChars := strings.Repeat(" ", tabs)
	attrs := reflect.ValueOf(o).Elem()
	for i := 0; i < attrs.NumField(); i++ {
		name := attrs.Type().Field(i).Name
		value := attrs.Field(i).Interface()
		color.Yellow("%s%s -> %s \n", tabChars, name, value)
	}
}
