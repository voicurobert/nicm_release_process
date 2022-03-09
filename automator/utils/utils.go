package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/sfreiberg/simplessh"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	skipRemovedDirs = []string{
		"release_patches",
		"dynamic_patches",
		"config",
		"nicm_dts_image",
		"nicm_register_web_services",
		"run",
		"run5"}
)

func DeleteFiles(args ...interface{}) error {
	rootDir := args[0].(string)
	fileExtension := args[1].(string)
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			for _, dirName := range skipRemovedDirs {
				if dirName == info.Name() {
					return filepath.SkipDir
				}
			}
		} else {
			if strings.HasSuffix(info.Name(), fileExtension) {
				_ = os.Remove(path)
			}
		}
		return nil
	})
}

func ExecuteGitPull(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", "git pull")
}

func DeleteJars(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", path+"delete_compiled_nicm.bat")
}

func CompileJars(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", path+"start_compile_all.bat")
}

func BuildJars(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	err = executeCommand("cmd.exe", "/C", path+"delete_compiled_nicm.bat")
	if err != nil {
		return err
	}
	err = executeCommand("cmd.exe", "/C", path+"start_compile_all.bat")
	if err != nil {
		return err
	}
	return nil
}

func SetScheduledTaskStatus(args ...interface{}) error {
	taskCommand := args[0].(string)
	disableTaskArgs := []string{"-NoProfile", "-NonInteractive"}
	disableTaskArgs = append(disableTaskArgs, taskCommand)
	disableTaskArgs = append(disableTaskArgs, "-TaskPath", "\"\\NICM\\\"", "-TaskName", "\"Test\"")
	return runPowerShellCommand(disableTaskArgs...)
}

func RunPowerShellScript(args ...interface{}) error {
	path := args[0].(string)
	return runPowerShellCommand([]string{path}...)
}

func runPowerShellCommand(args ...string) error {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		return err
	}
	return executeCommand(ps, args...)
}

func executeCommand(command string, args ...string) error {
	fmt.Println(args)
	cmd := exec.Command(command, args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			color.Yellow("\t > %s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		return err
	}
	return nil
}

func CreateArchive(args ...interface{}) error {
	path := args[0].(string)
	outName := args[1].(string)
	dirsToArchive := args[2].([]string)
	var skippedDirsFromArchive = make([]string, 0)
	if len(args) > 3 {
		skippedDirsFromArchive = args[3].([]string)
	}

	err := os.Chdir(path)
	if err != nil {
		color.Red("error changing dir %s: \n", err)
		return err
	}
	outFile, err := os.Create(path + outName)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	archive := zip.NewWriter(outFile)

	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {

		}
	}(archive)

	for _, path := range dirsToArchive {
		err := addFiles(archive, path, skippedDirsFromArchive)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFiles(w *zip.Writer, rootDir string, dirsToSkip []string) error {
	info, err := os.Stat(rootDir)
	if err != nil {
		return err
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
				return filepath.SkipDir
			}
			if strings.Contains(info.Name(), ".git") {
				return filepath.SkipDir
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

			}
		}(file)
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func isSkippedDir(dirName string, dirsToSkip []string) bool {
	if len(dirsToSkip) == 0 {
		return false
	}
	for _, dir := range dirsToSkip {
		if strings.Contains(dirName, dir) {
			return true
		}
	}
	return false
}

const (
	host     = "172.16.10.207"
	username = "laur"
	password = ""
)

func MoveArchive(args ...interface{}) error {
	sourcePath := args[0].(string)
	destPath := args[1].(string)

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()
	err = client.Upload(sourcePath, destPath)
	if err != nil {
		return err
	}

	return nil
}

func RenameArchive(args ...interface{}) error {
	currentName := args[0].(string)
	newName := args[1].(string)
	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Exec(fmt.Sprintf("mv %s %s", currentName, newName))

	if err != nil {
		return err
	}

	return nil
}

func DeleteOldArchive(args ...interface{}) error {
	name := args[0].(string)

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Exec(fmt.Sprintf("rm %s", name))

	if err != nil {
		return err
	}

	return nil
}

func Unzip(args ...interface{}) error {
	zipName := args[0].(string)

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Exec(fmt.Sprintf("unzip %s", zipName))

	if err != nil {
		return err
	}

	return nil
}

func BuildImage(args ...interface{}) error {
	shellBatPath := args[0].(string)

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Exec(fmt.Sprintf("run .sh from %s ", shellBatPath))

	if err != nil {
		return err
	}

	return nil
}
