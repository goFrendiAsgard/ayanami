package projectgenerator

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var mainCode = `package main
function main() {
	fmt.Println("nanana")
}
`

// ProjectGenerator configuration of generateProject
type ProjectGenerator struct {
	ProjectPath    string
	RepoName       string
	SourceCodePath string
	DeployablePath string
	GeneratorPath  string
}

// NewProjectGenerator create new project generator
func NewProjectGenerator(dirName, projectName, repoName string) (ProjectGenerator, error) {
	absDirPath, err := filepath.Abs(dirName)
	if err != nil {
		return ProjectGenerator{}, err
	}
	projectPath := filepath.Join(absDirPath, projectName)
	sourceCodePath := filepath.Join(projectPath, "sourcecode")
	deployablePath := filepath.Join(projectPath, "deployable")
	generatorPath := filepath.Join(projectPath, "generator")
	projectGenerator := ProjectGenerator{
		ProjectPath:    projectPath,
		SourceCodePath: sourceCodePath,
		DeployablePath: deployablePath,
		GeneratorPath:  generatorPath,
	}
	return projectGenerator, err
}

// Generate generating project skeleton
func (p ProjectGenerator) Generate() error {
	log.Println("Generate...")
	// create all directory
	err := os.MkdirAll(p.DeployablePath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(p.GeneratorPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(p.SourceCodePath, os.ModePerm)
	if err != nil {
		return err
	}
	// create `composition/main.go`
	mainPath := filepath.Join(p.SourceCodePath, "main.go")
	err = p.writeFile(mainPath, mainCode)
	if err != nil {
		return err
	}
	// run go mod init
	outByte, err := exec.Command("/bin/sh").Output()
	if err != nil {
		return err
	}
	log.Printf("[INFO] Init go module %s", string(outByte))
	return nil
}

func (p ProjectGenerator) writeFile(dstPath, content string) error {
	os.MkdirAll(dstPath, os.ModePerm)
	os.Remove(dstPath)
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return nil
}
