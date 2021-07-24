package action

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Action interface {
	DeleteFiles(...interface{}) error
}

type action struct {
}

func New() Action {
	return &action{}
}

func (a *action) DeleteFiles(args ...interface{}) error {
	rootDir := args[0].(string)
	fileExtension := args[1].(string)
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, "release_patches") || strings.Contains(path, "dynamic_patches") {
				return nil
			}
			ok, _ := regexp.MatchString(fileExtension, info.Name())
			if ok {
				err := os.Remove(path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
