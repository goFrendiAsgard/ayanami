package projectgenerator

import (
	"bytes"
	cp "github.com/otiai10/copy"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ProjectGenerator configuration of generateProject
type ProjectGenerator struct {
	ProjectSrcTemplatePath   string
	ProjectSrcGenPath        string
	ProjectPath              string
	RepoName                 string
	SourceCodePath           string
	DeployablePath           string
	GeneratorPath            string
	ProjectGeneratorTemplate *template.Template
}

// NewProjectGenerator create new project generator
func NewProjectGenerator(dirName, projectName, repoName, templatePath, genPath string) (ProjectGenerator, error) {
	projectGenerator := ProjectGenerator{}
	// get absolute dirPath of dirName
	absDirPath, err := filepath.Abs(dirName)
	if err != nil {
		return projectGenerator, err
	}
	// initiate and load template
	projectGeneratorTemplatePattern := filepath.Join(templatePath, "projectgenerator", "*")
	projectSrcTemplatePath := filepath.Join(templatePath, "project")
	projectGeneratorTemplate, err := template.ParseGlob(projectGeneratorTemplatePattern)
	if err != nil {
		return projectGenerator, err
	}
	// define directories
	projectPath := filepath.Join(absDirPath, projectName)
	sourceCodePath := filepath.Join(projectPath, "sourcecode")
	deployablePath := filepath.Join(projectPath, "deployable")
	generatorPath := filepath.Join(projectPath, "generator")
	// creat projectGenerator
	projectGenerator = ProjectGenerator{
		ProjectSrcTemplatePath:   projectSrcTemplatePath,
		ProjectSrcGenPath:        genPath,
		ProjectPath:              projectPath,
		RepoName:                 repoName,
		SourceCodePath:           sourceCodePath,
		DeployablePath:           deployablePath,
		GeneratorPath:            generatorPath,
		ProjectGeneratorTemplate: projectGeneratorTemplate,
	}
	return projectGenerator, err
}

// Generate generating project skeleton
func (p ProjectGenerator) Generate() error {
	log.Println("[INFO] Generate...")
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
	mainDstPath := filepath.Join(p.GeneratorPath, "main.go")
	err = p.Write(mainDstPath, "main.go", p)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", mainDstPath)
	// create `generator/go.mod`
	modDstPath := filepath.Join(p.GeneratorPath, "go.mod")
	err = p.Write(modDstPath, "go.mod", p)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", modDstPath)
	// create `generator/templates`
	templateDstPath := filepath.Join(p.GeneratorPath, "templates")
	templateSrcPath := p.ProjectSrcTemplatePath
	err = cp.Copy(templateSrcPath, templateDstPath)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", templateDstPath)
	// create `generator/templates`
	genDstPath := filepath.Join(p.GeneratorPath, "gen")
	genSrcPath := p.ProjectSrcGenPath
	err = cp.Copy(genSrcPath, genDstPath)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", genDstPath)
	return nil
}

// Write write using template
func (p ProjectGenerator) Write(dstPath, templateName string, data interface{}) error {
	buff := new(bytes.Buffer)
	p.ProjectGeneratorTemplate.ExecuteTemplate(buff, templateName, data)
	content := buff.String()
	content = strings.Trim(content, "\n")
	return p.WriteFile(dstPath, content)
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
