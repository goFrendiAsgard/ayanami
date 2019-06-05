package generator

import (
	"io"
	"os"
	"path/filepath"
)

// Resource Resource containing common methods to generate files
type Resource struct {
	cwd            string
	sourceCodePath string
	deployablePath string
}

// Copy copy resource from srcPath/src to dstPath/dst
func (r Resource) Copy(src, dst string) error {
	// read from src
	src = filepath.Join(r.sourceCodePath, src)
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	// create dst's parent directory if not exists
	dst = filepath.Join(r.deployablePath, dst)
	mkdirAll(dst)
	// write to dst
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	// copy
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// WriteString write content to dstPath/dst
func (r Resource) WriteString(dst string, content string) error {
	dst = filepath.Join(r.deployablePath, dst)
	mkdirAll(dst)
	os.Remove(dst)
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	return nil
}

// NewResource create new resource
func NewResource() (Resource, error) {
	resource := Resource{}
	cwd, err := os.Getwd()
	if err != nil {
		return resource, err
	}
	pwd := filepath.Dir(cwd)
	resource.cwd = cwd
	resource.sourceCodePath = filepath.Join(pwd, "sourcecode")
	resource.deployablePath = filepath.Join(pwd, "deployable")
	return resource, err
}

func mkdirAll(dst string) {
	dstParent := filepath.Dir(dst)
	os.MkdirAll(dstParent, os.ModePerm)
}
