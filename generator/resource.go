package generator

import (
	"bytes"
	cp "github.com/otiai10/copy"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Resource Resource containing common methods to generate files
type Resource struct {
	SourceCodePath string
	DeployablePath string
	Template       *template.Template
}

// NewResource create new resource
func NewResource(sourceCodePath, deployablePath, templatePath string) (Resource, error) {
	projectTemplatePattern := filepath.Join(templatePath, "*")
	projectTemplate := template.Must(template.ParseGlob(projectTemplatePattern))
	resource := Resource{
		SourceCodePath: sourceCodePath,
		DeployablePath: deployablePath,
		Template:       projectTemplate,
	}
	return resource, nil
}

// NewResourceByProjectPath create new resource using
func NewResourceByProjectPath(projectPath string) (Resource, error) {
	sourceCodePath := filepath.Join(projectPath, "sourcecode")
	deployablePath := filepath.Join(projectPath, "deployable")
	templatePath := filepath.Join(projectPath, "generator", "templates")
	return NewResource(sourceCodePath, deployablePath, templatePath)
}

// Copy src to dst
func (r Resource) Copy(src, dst string) error {
	// make sure parent directory exists
	dstParent := filepath.Dir(dst)
	err := os.MkdirAll(dstParent, os.ModePerm)
	if err != nil {
		return err
	}
	// copy
	return cp.Copy(src, dst)
}

// CopySourceToDeployable copy sourcecode/src to deployable/dst
func (r Resource) CopySourceToDeployable(src, dst string) error {
	// read from src
	src = filepath.Join(r.SourceCodePath, src)
	// create dst's parent directory if not exists
	dst = filepath.Join(r.DeployablePath, dst)
	return r.Copy(src, dst)
}

// WriteDeployable write to sourcecode/dstPath
func (r Resource) WriteDeployable(dstPath, templateName string, data interface{}) error {
	dstPath = filepath.Join(r.DeployablePath, dstPath)
	return r.Write(dstPath, templateName, data)
}

// WriteSource write to sourcecode/dstPath
func (r Resource) WriteSource(dstPath, templateName string, data interface{}) error {
	dstPath = filepath.Join(r.SourceCodePath, dstPath)
	return r.Write(dstPath, templateName, data)
}

// Write write using template
func (r Resource) Write(dstPath, templateName string, data interface{}) error {
	buff := new(bytes.Buffer)
	r.Template.ExecuteTemplate(buff, templateName, data)
	content := buff.String()
	content = strings.Trim(content, "\n")
	return r.WriteFile(dstPath, content)
}

// WriteFile write content to filePath
func (r Resource) WriteFile(filePath string, content string) error {
	// make sure parent directory exists
	fileParentPath := filepath.Dir(filePath)
	err := os.MkdirAll(fileParentPath, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return err
}
