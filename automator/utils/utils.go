package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DeleteFiles(args ...interface{}) error {
	rootDir := args[0].(string)
	fileExtension := args[1].(string)
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.Contains(path, "release_patches") ||
				strings.Contains(path, "dynamic_patches") ||
				strings.Contains(path, "build_nicm_linux") {
				return nil
			}
			if strings.HasSuffix(info.Name(), fileExtension) {
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
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", "git pull")
}

func BuildImages(args ...interface{}) error {
	path := args[0].(string)
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	return executeCommand("cmd.exe", "/C", path+"build_all.bat")
}

func SetScheduledTaskStatus(args ...interface{}) error {
	disableTask := args[0].(string)
	disableTaskArgs := []string{"-NoProfile", "-NonInteractive"}
	var taskCommand string
	if disableTask == "true" {
		taskCommand = "Disable-ScheduledTask"
	} else {
		taskCommand = "Enable-ScheduledTask"
	}
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
	imageNames := args[1].([]string)
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
	var dirsToSkip = make([]string, 0)
	if len(args) > 3 {
		dirsToSkip = args[3].([]string)
	}

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

	for _, path := range dirsToArchive {
		err := addFiles(archive, path, dirsToSkip)
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

func MoveArchive(args ...interface{}) error {
	sourcePath := args[0].(string)
	destPath := args[1].(string)
	archiveName := args[2].(string)
	//cmdArgs := []string{"-p", "T#ink2!", "scp", "-r", sourcePath, "laur", "@", "172.16.10.207", ":", destPath + "//" + archiveName}
	arg := fmt.Sprintf("-p %s %s %s@%s:%s", "T#ink2!", sourcePath+archiveName, "laur", "172.16.10.207", destPath+archiveName)
	fmt.Println(arg)
	cmdArgs := []string{arg}
	return executeCommand("ssh", cmdArgs...)
}
