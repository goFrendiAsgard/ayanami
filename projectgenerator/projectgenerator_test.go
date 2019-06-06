package projectgenerator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProjectGenerator(t *testing.T) {
	_, currentDirPath, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("fail to locate directory")
	}
	ayanamiDirPath := filepath.Dir(filepath.Dir(currentDirPath))
	projectParentDirPath := filepath.Join(ayanamiDirPath, ".test-projectgenerator")
	templatePath := filepath.Join(ayanamiDirPath, "templates")
	genPath := filepath.Join(ayanamiDirPath, "gen")
	// create generator
	generator, err := NewProjectGenerator(projectParentDirPath, "evangelion", "github.com/nerv/evangelion", templatePath, genPath)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	// generate
	err = generator.Generate()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	defer os.RemoveAll(projectParentDirPath)

	// check deployable
	dirPath := filepath.Join(projectParentDirPath, "evangelion", "deployable")
	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		t.Errorf("%s is not exists", dirPath)
	}
	// check sourcecode
	dirPath = filepath.Join(projectParentDirPath, "evangelion", "sourcecode")
	if stat, err := os.Stat(dirPath); err != nil || !stat.IsDir() {
		t.Errorf("%s is not exists", dirPath)
	}

	// check generator/go.mod content
	gomodFile := filepath.Join(projectParentDirPath, "evangelion", "generator", "go.mod")
	gomodByte, err := ioutil.ReadFile(gomodFile)
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	expectedGoModContent := "module github.com/nerv/evangelion"
	actualGoModContent := string(gomodByte)
	if expectedGoModContent != actualGoModContent {
		t.Errorf("expected `%s`, get `%s`", expectedGoModContent, actualGoModContent)
	}

	// check generator/main.go existance
	if _, err := os.Stat(filepath.Join(projectParentDirPath, "evangelion", "generator", "main.go")); err != nil {
		t.Errorf("Get error: %s", err)
	}

	// check generator/templates existance
	if _, err := os.Stat(filepath.Join(projectParentDirPath, "evangelion", "generator", "templates")); err != nil {
		t.Errorf("Get error: %s", err)
	}

	// check generator/gen existance
	if _, err := os.Stat(filepath.Join(projectParentDirPath, "evangelion", "generator", "gen")); err != nil {
		t.Errorf("Get error: %s", err)
	}

}
