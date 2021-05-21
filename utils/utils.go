package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func DeleteFiles(args ...interface{}) error {
	rootDir := args[0].(string)
	fileExtension := args[1].(string)
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
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

func ExecuteGitPull(args ...interface{}) error {
	rootDir := args[0].(string)
	repo, err := git.PlainOpen(rootDir)
	if err != nil {
		return err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = wt.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: "refs/heads/master",
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
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", path+"build_all.bat")
}

func SetTaskStatus(args ...interface{}) error {
	scriptPath := args[0].(string)
	disable := args[1].(string)
	return executeCommand("cmd.exe", "/C", scriptPath, disable)
}

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		if err2 != nil {
			return err2
		}
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		if err2 != nil {
			return err2
		}
		return err
	}
	err = cmd.Wait()
	if err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		if err2 != nil {
			return err2
		}
		return err
	}
	return nil
}

func SetWritableAccess(args ...interface{}) error {
	fullPath := args[0].(string)
	imageNames := args[2].([]string)
	return filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			for _, imageName := range imageNames {
				if strings.Contains(path, imageName) {
					err := os.Chmod(path, 0222)
					if err != nil {
						fmt.Printf("error setting writable access %s \n", err.Error())
						return err
					}
				}
			}
		}
		return nil
	})
}

func CreateArchive(args ...interface{}) error {
	path := args[0].(string)
	outName := args[1].(string)
	dirsToArchive := args[2].([]string)
	dirsToSkip := args[3].([]string)
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
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			fmt.Println("error when trying to close out file")
		}
	}(outFile)
	archive := zip.NewWriter(outFile)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			fmt.Println("error when trying to close archiver")
		}
	}(archive)

	for _, path := range dirsToArchive {
		addFiles(archive, path, dirsToSkip)
	}
	return nil
}

func addFiles(w *zip.Writer, rootDir string, dirsToSkip []string) {
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
			if isSkippedDir(info.Name(), dirsToSkip) {
				return nil
			}
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
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("error when trying to close file")
			}
		}(file)
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func isSkippedDir(dirName string, dirsToSkip []string) bool {
	if len(dirsToSkip) == 0 {
		return false
	}
	for _, dir := range dirsToSkip {
		if dirName == dir {
			return true
		}
	}
	return false
}
