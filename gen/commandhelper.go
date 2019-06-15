package gen

import (
	"github.com/state-alchemists/ayanami/generator"
	"os/exec"
	"path/filepath"
)

// GitInit run gitInit in a directory
func GitInit(io generator.IOHelper, dirPath string) error {
	gitPath := filepath.Join(dirPath, ".git")
	if !io.IsExists(gitPath) {
		shellCmd := exec.Command("git", "init")
		shellCmd.Dir = dirPath
		err := shellCmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// GoFmt run gofmt in a directory
func GoFmt(dirPath string) error {
	shellCmd := exec.Command("gofmt", "-w", "-s", ".")
	shellCmd.Dir = dirPath
	err := shellCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
