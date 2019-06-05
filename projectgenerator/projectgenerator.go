package projectgenerator

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var mainCode = `package main

import(
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
)

// Generator our main generator
var Generator = generator.NewGenerator()

func init() {
	// do something here
}

func main() {
	fmt.Println("nanana")
}
`
var modContent = `module {{.RepoName}}`
var modTemplate = template.Must(template.New("mod").Parse(modContent))

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
		RepoName:       repoName,
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
	// create deployable directory
	err := os.MkdirAll(p.DeployablePath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", p.DeployablePath)
	// create generator directory
	err = os.MkdirAll(p.GeneratorPath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", p.GeneratorPath)
	// create sourcode directory
	err = os.MkdirAll(p.SourceCodePath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", p.SourceCodePath)
	// create `generator/main.go`
	mainPath := filepath.Join(p.GeneratorPath, "main.go")
	err = p.WriteFile(mainPath, mainCode)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", mainPath)
	// create `composition/go.mod`
	modPath := filepath.Join(p.GeneratorPath, "go.mod")
	err = p.WriteTemplate(modPath, modTemplate, p)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", modPath)
	return nil
}

// WriteTemplate write using template
func (p ProjectGenerator) WriteTemplate(dstPath string, template *template.Template, data interface{}) error {
	buff := new(bytes.Buffer)
	template.Execute(buff, data)
	return p.WriteFile(dstPath, buff.String())
}

// WriteFile write content to file
func (p ProjectGenerator) WriteFile(dstPath, content string) error {
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
