package utils

import (
	"archive/zip"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DeleteFiles(args ...interface{}) error {
	return filepath.Walk(args[0].(string), func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.Contains(path, ".magikc") {
			return os.Remove(path)
		}
		return nil
	})
}

func ExecuteGitPull(args ...interface{}) error {
	repo, err := git.PlainOpen(args[0].(string))
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "ds.nicm",
			Password: "aergistal",
		},
		Progress: os.Stdout,
	})
	if err != nil {
		if strings.Contains(err.Error(), "up-to-date") {
			return nil
		}
		return err
	}
	return nil
}

func BuildImages(args ...interface{}) error {
	wk, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(args[0].(string))
	if err != nil {
		return err
	}
	cmd := exec.Command("cmd.exe", "/C", "ant", args[1].(string))
	//cmd.Env = os.Environ()
	//cmd.Env = append(cmd.Env, "MY_VAR=some_value")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		log.Print(string(bytes))
		fmt.Println(string(bytes))
		return err
	}
	log.Print(string(bytes))
	fmt.Println(string(bytes))
	err = os.Chdir(wk)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func SetWritableAccess(args ...interface{}) error {
	fullPath := args[0].(string)
	return filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, args[1].(string)) || strings.Contains(path, args[2].(string)) {
				err := os.Chmod(path, 0222)
				if err != nil {
					fmt.Printf("error setting writable access %s \n", err.Error())
					return err
				}
			}
		}
		return nil
	})
}

func CreateArchive(args ...interface{}) error {
	path := args[0].(string)
	outName := args[1].(string)
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("error changing dir %s: \n", err)
		return err
	}
	outFile, err := os.Create(path + outName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer outFile.Close()
	archive := zip.NewWriter(outFile)
	defer archive.Close()

	for _, path := range args[2:] {
		addFiles(archive, path.(string))
	}
	return nil
}

func addFiles(w *zip.Writer, rootDir string) {
	info, err := os.Stat(rootDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(rootDir)
	}

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, rootDir))
		}
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := w.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
