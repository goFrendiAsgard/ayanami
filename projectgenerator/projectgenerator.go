package projectgenerator

import (
	"bytes"
	"fmt"
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
	ExampleType              string
	ProjectGeneratorTemplate *template.Template
}

// Generate generating project skeleton
func (pg ProjectGenerator) Generate() error {
	log.Println("[INFO] Generate...")
	// create deployable directory
	log.Printf("[INFO] Create %s", pg.DeployablePath)
	err := os.MkdirAll(pg.DeployablePath, os.ModePerm)
	if err != nil {
		return err
	}
	// create sourcode directory
	log.Printf("[INFO] Create %s", pg.SourceCodePath)
	err = os.MkdirAll(pg.SourceCodePath, os.ModePerm)
	if err != nil {
		return err
	}
	// create generator directory
	log.Printf("[INFO] Create %s", pg.GeneratorPath)
	err = os.MkdirAll(pg.GeneratorPath, os.ModePerm)
	if err != nil {
		return err
	}
	// create `generator/templates`
	templateDstPath := filepath.Join(pg.GeneratorPath, "templates")
	templateSrcPath := pg.ProjectSrcTemplatePath
	log.Printf("[INFO] Create %s", templateDstPath)
	err = cp.Copy(templateSrcPath, templateDstPath)
	if err != nil {
		return err
	}
	// create `generator/gen`
	genDstPath := filepath.Join(pg.GeneratorPath, "gen")
	genSrcPath := pg.ProjectSrcGenPath
	log.Printf("[INFO] Create %s", genDstPath)
	err = cp.Copy(genSrcPath, genDstPath)
	if err != nil {
		return err
	}
	// create `generator/<whatever>` from non exampleType
	for _, templateName := range pg.getExampleNonTypeTemplates() {
		dstFileName := templateName
		dstPath := filepath.Join(pg.GeneratorPath, dstFileName)
		log.Printf("[INFO] Create %s", dstPath)
		err = pg.Write(dstPath, templateName, pg)
		if err != nil {
			return err
		}
	}
	// create `generator/<whatever>` from exampleType
	for _, templateName := range pg.getExampleTypeTemplates() {
		templateParts := strings.Split(templateName, ".")
		dstFileName := strings.Join(templateParts[2:], ".")
		dstPath := filepath.Join(pg.GeneratorPath, dstFileName)
		log.Printf("[INFO] Create %s", dstPath)
		err = pg.Write(dstPath, templateName, pg)
		if err != nil {
			return err
		}
	}
	return nil
}

// Write write using template
func (pg ProjectGenerator) Write(dstPath, templateName string, data interface{}) error {
	buff := new(bytes.Buffer)
	err := pg.ProjectGeneratorTemplate.ExecuteTemplate(buff, templateName, data)
	if err != nil {
		return err
	}
	content := buff.String()
	content = strings.Trim(content, "\n")
	return pg.WriteFile(dstPath, content)
}

// WriteFile write content to file
func (pg ProjectGenerator) WriteFile(dstPath, content string) error {
	err := os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.Remove(dstPath)
	if err != nil {
		return err
	}
	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("[ERROR] Cannot close `%s`: %s", dstPath, err)
		}
	}()
	_, err = f.WriteString(content)
	return nil
}

func (pg ProjectGenerator) getExampleNonTypeTemplates() []string {
	definedTemplateString := pg.ProjectGeneratorTemplate.DefinedTemplates()
	definedTemplateString = definedTemplateString[len("; defined templates are: "):]
	definedTemplates := strings.Split(definedTemplateString, ", ")
	var exampleTemplates []string
	for index := range definedTemplates {
		templateName := strings.Trim(definedTemplates[index], " ")
		templateName = strings.Trim(definedTemplates[index], `"`)
		if templateName[len(templateName)-5:] == ".tmpl" {
			continue
		}
		if strings.Index(templateName, "example.") != 0 {
			exampleTemplates = append(exampleTemplates, templateName)
		}
	}
	return exampleTemplates
}

func (pg ProjectGenerator) getExampleTypeTemplates() []string {
	definedTemplateString := pg.ProjectGeneratorTemplate.DefinedTemplates()
	definedTemplateString = definedTemplateString[len("; defined templates are: "):]
	definedTemplates := strings.Split(definedTemplateString, ", ")
	var exampleTemplates []string
	for index := range definedTemplates {
		templateName := strings.Trim(definedTemplates[index], " ")
		templateName = strings.Trim(definedTemplates[index], `"`)
		if templateName[len(templateName)-5:] == ".tmpl" {
			continue
		}
		if strings.Index(templateName, fmt.Sprintf("example.%s.", pg.ExampleType)) == 0 {
			exampleTemplates = append(exampleTemplates, templateName)
		}
	}
	return exampleTemplates
}

// NewProjectGenerator create new project generator
func NewProjectGenerator(dirName, projectName, repoName, templatePath, genPath string, exampleType string) (ProjectGenerator, error) {
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
		ExampleType:              exampleType,
	}
	return projectGenerator, err
}
