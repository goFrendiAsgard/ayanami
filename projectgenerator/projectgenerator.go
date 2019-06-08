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

// Generate generating project skeleton
func (pg ProjectGenerator) Generate() error {
	log.Println("[INFO] Generate...")
	// create deployable directory
	err := os.MkdirAll(pg.DeployablePath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", pg.DeployablePath)
	// create generator directory
	err = os.MkdirAll(pg.GeneratorPath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", pg.GeneratorPath)
	// create sourcode directory
	err = os.MkdirAll(pg.SourceCodePath, os.ModePerm)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", pg.SourceCodePath)
	// create `generator/main.go`
	mainDstPath := filepath.Join(pg.GeneratorPath, "main.go")
	err = pg.Write(mainDstPath, "main.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", mainDstPath)
	// create `generator/go.mod`
	modDstPath := filepath.Join(pg.GeneratorPath, "go.mod")
	err = pg.Write(modDstPath, "go.mod", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", modDstPath)
	// create `generator/templates`
	templateDstPath := filepath.Join(pg.GeneratorPath, "templates")
	templateSrcPath := pg.ProjectSrcTemplatePath
	err = cp.Copy(templateSrcPath, templateDstPath)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", templateDstPath)
	// create `generator/templates`
	genDstPath := filepath.Join(pg.GeneratorPath, "gen")
	genSrcPath := pg.ProjectSrcGenPath
	err = cp.Copy(genSrcPath, genDstPath)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", genDstPath)
	// create `generator/gateway.go`
	gatewayDstPath := filepath.Join(pg.GeneratorPath, "gateway.go")
	err = pg.Write(gatewayDstPath, "example.gateway.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", gatewayDstPath)
	// create `generator/gomonolith.go`
	gomonolithDstPath := filepath.Join(pg.GeneratorPath, "monolith.go")
	err = pg.Write(gomonolithDstPath, "example.gomonolith.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", gomonolithDstPath)
	// create `generator/cmdservice.go`
	cmdServiceDstPath := filepath.Join(pg.GeneratorPath, "cmdservice.go")
	err = pg.Write(cmdServiceDstPath, "example.cmd.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", cmdServiceDstPath)
	// create `generator/htmlservice.go`
	htmlServiceDstPath := filepath.Join(pg.GeneratorPath, "htmlservice.go")
	err = pg.Write(htmlServiceDstPath, "example.gosrvc.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", htmlServiceDstPath)
	// create `generator/flowbanner.go`
	flowBannerServiceDstPath := filepath.Join(pg.GeneratorPath, "flowbanner.go")
	err = pg.Write(flowBannerServiceDstPath, "example.flow.banner.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", flowBannerServiceDstPath)
	// create `generator/flowroot.go`
	flowRootServiceDstPath := filepath.Join(pg.GeneratorPath, "flowroot.go")
	err = pg.Write(flowRootServiceDstPath, "example.flow.root.go", pg)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Create %s", flowRootServiceDstPath)
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
	defer f.Close()
	_, err = f.WriteString(content)
	return nil
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
