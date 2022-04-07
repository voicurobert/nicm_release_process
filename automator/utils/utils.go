package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/fatih/color"
	"github.com/sfreiberg/simplessh"
	"github.com/voicurobert/nicm_release_process/automator/process/options"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	skipRemovedDirs = []string{
		"release_patches",
		"dynamic_patches",
		"hotfixes",
		"config",
		"nicm_dts_image",
		"nicm_register_web_services",
		"run",
		"run5",
	}

	skipDirs = []string{
		".git",
		"build_nicm_linux",
		"nicm_munit",
		"nicm_nig",
		"nicm_nig_resources",
		"nicm_nig_web",
		"nicm_probe",
		"nicm_upgrade_framework",
		"nicm_night_scripts",
		"nicm_nig",
	}
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

func BuildJars(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
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

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error creating StderrPipe for Cmd", err)
		return err
	}

	errScanner := bufio.NewScanner(cmdErr)
	go func() {
		for errScanner.Scan() {
			color.Red("\t > %s\n", errScanner.Text())
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
	defer outFile.Close()

	archive := zip.NewWriter(outFile)
	defer archive.Close()

	for _, subDir := range dirsToArchive {
		err := addFiles(archive, subDir, skipDirs)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFiles(w *zip.Writer, rootDir string, dirsToSkip []string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if isSkippedDir(info.Name(), dirsToSkip) {
				return filepath.SkipDir
			}
			if strings.Contains(info.Name(), ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		path = strings.TrimPrefix(path, rootDir+"/")

		if strings.HasSuffix(path, ".magik") {
			if !skipDirsFromMagikFiles(path) {
				return nil
			}
		}
		f, err := w.Create(path)
		if err != nil {
			return err
		}
		file, err := os.Open(path)

		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(f, file)
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

func skipDirsFromMagikFiles(path string) bool {
	names := []string{
		"nicm_products\\nicm\\config",
		"nicm_products\\nicm_build",
		"nicm_products\\nicm\\source\\release_patches",
		"nicm_products\\nicm\\hotfixes",
		"nicm_products\\rwee_extensions",
		"nicm\\dynamic_patches"}

	for _, name := range names {
		if strings.Contains(path, name) {
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

func RunRemoteCommands(args ...interface{}) error {

	screenName := args[0]

	var client *simplessh.Client
	var err error

	if client, err = simplessh.ConnectWithPassword(host, username, password); err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Exec(fmt.Sprintf("screen -DmS %s quit()", screenName))
	if err != nil {
		return err
	}

	_, err = client.Exec(fmt.Sprintf("screen -DmS %s start ....", screenName))
	if err != nil {
		return err
	}

	return nil
}

type ConfigMap map[string]map[string]string

func getConfigPath() string {
	dir, _ := os.Getwd()
	return path.Join(dir, "nicm_paths.config")
}

func GetConfig() ConfigMap {
	configFilePath := getConfigPath()
	config, err := configparser.NewConfigParserFromFile(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("cannot read config file %s, error: %s", configFilePath, err.Error()))
	}

	configMap := make(ConfigMap)

	for _, section := range config.Sections() {
		keyValue, _ := config.Items(section)
		configMap[section] = keyValue
	}

	return configMap
}

func SetOptionPaths(options *options.Options, cfgMap map[string]string) {
	if gitPath, ok := cfgMap["git_path"]; ok {
		options.SetGitPath(strings.TrimSpace(gitPath))
	}
	if workingPath, ok := cfgMap["working_path"]; ok {
		options.SetWorkingPath(strings.TrimSpace(workingPath))
	}
}
