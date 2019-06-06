package generator

import (
	"bytes"
	cp "github.com/otiai10/copy"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// IOHelper IOHelper containing common methods to generate files
type IOHelper struct {
	sourcePath string
	depPath    string
	template   *template.Template
}

// NewIOHelper create new IOHelper
func NewIOHelper(sourcePath, depPath, templatePath string) (IOHelper, error) {
	projectTemplatePattern := filepath.Join(templatePath, "*")
	projectTemplate := template.Must(template.ParseGlob(projectTemplatePattern))
	io := IOHelper{
		sourcePath: sourcePath,
		depPath:    depPath,
		template:   projectTemplate,
	}
	return io, nil
}

// NewIOHelperByProjectPath create new io by using projectPath
func NewIOHelperByProjectPath(projectPath string) (IOHelper, error) {
	sourceCodePath := filepath.Join(projectPath, "sourcecode")
	deployablePath := filepath.Join(projectPath, "deployable")
	templatePath := filepath.Join(projectPath, "generator", "templates")
	return NewIOHelper(sourceCodePath, deployablePath, templatePath)
}

// GetSourcePath get sourcePath
func (io *IOHelper) GetSourcePath() string {
	return io.sourcePath
}

// GetDepPath get depPath
func (io *IOHelper) GetDepPath() string {
	return io.depPath
}

// GetTemplate get template
func (io *IOHelper) GetTemplate() *template.Template {
	return io.template
}

// IsExists check whether filePath exists or not
func (io *IOHelper) IsExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// IsSourceExists check whether filePath exists or not
func (io *IOHelper) IsSourceExists(filePath string) bool {
	filePath = filepath.Join(io.sourcePath, filePath)
	return io.IsExists(filePath)
}

// IsDepExists check whether filePath exists or not
func (io *IOHelper) IsDepExists(filePath string) bool {
	filePath = filepath.Join(io.depPath, filePath)
	return io.IsExists(filePath)
}

// Copy src to dst
func (io *IOHelper) Copy(src, dst string) error {
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
func (io *IOHelper) CopySourceToDep(src, dst string) error {
	// read from src
	src = filepath.Join(io.sourcePath, src)
	// create dst's parent directory if not exists
	dst = filepath.Join(io.depPath, dst)
	return io.Copy(src, dst)
}

// WriteDep write to sourcecode/dstPath
func (io *IOHelper) WriteDep(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(io.depPath, filePath)
	return io.Write(filePath, templateName, data)
}

// WriteSource write to sourcecode/dstPath
func (io *IOHelper) WriteSource(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(io.sourcePath, filePath)
	return io.Write(filePath, templateName, data)
}

// Write write using template
func (io *IOHelper) Write(filePath, templateName string, data interface{}) error {
	buff := new(bytes.Buffer)
	err := io.template.ExecuteTemplate(buff, templateName, data)
	if err != nil {
		return err
	}
	content := buff.String()
	content = strings.Trim(content, "\n")
	return io.WriteFile(filePath, content)
}

// WriteFile write content to filePath
func (io *IOHelper) WriteFile(filePath string, content string) error {
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
