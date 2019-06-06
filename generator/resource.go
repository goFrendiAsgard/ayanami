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
	sourcePath string
	depPath    string
	template   *template.Template
}

// NewResource create new resource
func NewResource(sourcePath, depPath, templatePath string) (*Resource, error) {
	projectTemplatePattern := filepath.Join(templatePath, "*")
	projectTemplate := template.Must(template.ParseGlob(projectTemplatePattern))
	resource := Resource{
		sourcePath: sourcePath,
		depPath:    depPath,
		template:   projectTemplate,
	}
	return &resource, nil
}

// NewResourceByProjectPath create new resource using
func NewResourceByProjectPath(projectPath string) (*Resource, error) {
	sourceCodePath := filepath.Join(projectPath, "sourcecode")
	deployablePath := filepath.Join(projectPath, "deployable")
	templatePath := filepath.Join(projectPath, "generator", "templates")
	return NewResource(sourceCodePath, deployablePath, templatePath)
}

// GetSourcePath get sourcePath
func (r *Resource) GetSourcePath() string {
	return r.sourcePath
}

// GetDepPath get depPath
func (r *Resource) GetDepPath() string {
	return r.depPath
}

// GetTemplate get template
func (r *Resource) GetTemplate() *template.Template {
	return r.template
}

// IsExists check whether filePath exists or not
func (r *Resource) IsExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// IsSourceExists check whether filePath exists or not
func (r *Resource) IsSourceExists(filePath string) bool {
	filePath = filepath.Join(r.sourcePath, filePath)
	return r.IsExists(filePath)
}

// IsDepExists check whether filePath exists or not
func (r *Resource) IsDepExists(filePath string) bool {
	filePath = filepath.Join(r.depPath, filePath)
	return r.IsExists(filePath)
}

// Copy src to dst
func (r *Resource) Copy(src, dst string) error {
	// make sure parent directory exists
	dstParent := filepath.Dir(dst)
	err := os.MkdirAll(dstParent, os.ModePerm)
	if err != nil {
		return err
	}
	// copy
	return cp.Copy(src, dst)
}

// CopySourceToDep copy sourcecode/src to deployable/dst
func (r *Resource) CopySourceToDep(src, dst string) error {
	// read from src
	src = filepath.Join(r.sourcePath, src)
	// create dst's parent directory if not exists
	dst = filepath.Join(r.depPath, dst)
	return r.Copy(src, dst)
}

// WriteDep write to sourcecode/dstPath
func (r *Resource) WriteDep(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(r.depPath, filePath)
	return r.Write(filePath, templateName, data)
}

// WriteSource write to sourcecode/dstPath
func (r *Resource) WriteSource(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(r.sourcePath, filePath)
	return r.Write(filePath, templateName, data)
}

// Write write using template
func (r *Resource) Write(filePath, templateName string, data interface{}) error {
	buff := new(bytes.Buffer)
	err := r.template.ExecuteTemplate(buff, templateName, data)
	if err != nil {
		return err
	}
	content := buff.String()
	content = strings.Trim(content, "\n")
	return r.WriteFile(filePath, content)
}

// WriteFile write content to filePath
func (r *Resource) WriteFile(filePath string, content string) error {
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
