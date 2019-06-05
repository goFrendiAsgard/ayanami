package projectgenerator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectGenerator(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	// create generator
	dirName := filepath.Join(cwd, ".test")
	generator, err := NewProjectGenerator(dirName, "evangelion", "github.com/nerv/evangelion")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	// generate
	err = generator.Generate()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	// check deployable
	dirPath := filepath.Join(dirName, "evangelion", "deployable")
	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		t.Errorf("%s is not exists", dirPath)
	}
	// check sourcecode
	dirPath = filepath.Join(dirName, "evangelion", "sourcecode")
	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		t.Errorf("%s is not exists", dirPath)
	}
	// check generator/go.mod
	gomodFile := filepath.Join(dirName, "evangelion", "generator", "go.mod")
	gomodByte, err := ioutil.ReadFile(gomodFile)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expectedGoModContent := "module github.com/nerv/evangelion"
	actualGoModContent := string(gomodByte)
	if expectedGoModContent != actualGoModContent {
		t.Errorf("expected `%s`, get `%s`", expectedGoModContent, actualGoModContent)
	}
	// remove test dir
	os.RemoveAll(dirName)
}
