package generator

import (
	"bytes"
	cp "github.com/otiai10/copy"
	"log"
	"os"
	"os/exec"
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

// GitInitSource run git init for a directory in source path
func (io *IOHelper) GitInitSource(dirPath string) error {
	return io.GitInit(filepath.Join(io.sourcePath, dirPath))
}

// GitInitDep run git init for a directory in dep path
func (io *IOHelper) GitInitDep(dirPath string) error {
	return io.GitInit(filepath.Join(io.depPath, dirPath))
}

// GitInit run git init in a directory
func (io *IOHelper) GitInit(dirPath string) error {
	err := io.MkdirAll(dirPath)
	if err != nil {
		return err
	}
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

// GoFmtSource run git init for a directory in source path
func (io *IOHelper) GoFmtSource(dirPath string) error {
	return io.GoFmt(filepath.Join(io.sourcePath, dirPath))
}

// GoFmtDep run git init for a directory in dep path
func (io *IOHelper) GoFmtDep(dirPath string) error {
	return io.GoFmt(filepath.Join(io.sourcePath, dirPath))
}

// GoFmt run gofmt in a directory
func (io *IOHelper) GoFmt(dirPath string) error {
	err := io.MkdirAll(dirPath)
	if err != nil {
		return err
	}
	shellCmd := exec.Command("gofmt", "-w", "-s", ".")
	shellCmd.Dir = dirPath
	err = shellCmd.Run()
	if err != nil {
		return err
	}
	return nil
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
	err := io.MkdirAll(dstParent)
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

// WriteDeps write everything to dep
func (io *IOHelper) WriteDeps(dirPath string, templateNames []string, data interface{}) error {
	// write common things
	for _, templateName := range templateNames {
		log.Printf("[INFO] Create %s", templateName)
		filePath := filepath.Join(dirPath, templateName)
		err := io.WriteDep(filePath, templateName, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteDep write to sourcecode/dstPath
func (io *IOHelper) WriteDep(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(io.depPath, filePath)
	return io.Write(filePath, templateName, data)
}

// WriteSources write everything to sources
func (io *IOHelper) WriteSources(dirPath string, templateNames []string, data interface{}) error {
	// write common things
	for _, templateName := range templateNames {
		log.Printf("[INFO] Create %s", templateName)
		filePath := filepath.Join(dirPath, templateName)
		err := io.WriteSource(filePath, templateName, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteSource write to sourcecode/dstPath
func (io *IOHelper) WriteSource(filePath, templateName string, data interface{}) error {
	filePath = filepath.Join(io.sourcePath, filePath)
	return io.Write(filePath, templateName, data)
}

// Write write using template
func (io *IOHelper) Write(filePath, templateName string, data interface{}) error {
	content, err := io.GetParsedTemplate(templateName, data)
	if err != nil {
		return err
	}
	return io.WriteFile(filePath, content)
}

// GetParsedTemplate get parsed template
func (io *IOHelper) GetParsedTemplate(templateName string, data interface{}) (string, error) {
	buff := new(bytes.Buffer)
	err := io.template.ExecuteTemplate(buff, templateName, data)
	if err != nil {
		return "", err
	}
	content := buff.String()
	content = strings.Trim(content, "\n")
	return content, nil
}

// MkdirAll create directory and it's parent if necessary
func (io *IOHelper) MkdirAll(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

// WriteFile write content to filePath
func (io *IOHelper) WriteFile(filePath string, content string) error {
	// make sure parent directory exists
	fileParentPath := filepath.Dir(filePath)
	err := io.MkdirAll(fileParentPath)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close file `%s`: %s", filePath, err)
		}
	}()
	_, err = f.WriteString(content)
	return err
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
