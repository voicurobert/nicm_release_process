package options

import (
	"fmt"
	"reflect"
	"strings"
)

type Options struct {
	WorkingPath string
	GitPath     string
	BuildPath   string
	AntCommand  string
	ImagesPath  string
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

func (o *Options) Print() {
	attrs := reflect.ValueOf(o).Elem()
	for i := 0; i < attrs.NumField(); i++ {
		name := attrs.Type().Field(i).Name
		value := attrs.Field(i).Interface()
		fmt.Printf("Option: %s --- Value: %s \n", name, value)
	}
}
