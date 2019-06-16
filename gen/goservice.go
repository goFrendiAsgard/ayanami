package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
	"strings"
)

// ExposedGoServiceConfig exposed ready flowConfig
type ExposedGoServiceConfig struct {
	MainFunctionName string
	ServiceName      string
	RepoName         string
	Packages         []string
	Functions        map[string]ExposedFunction
}

// GoServiceConfig configuration to generate GoService
type GoServiceConfig struct {
	ServiceName string
	RepoName    string
	Functions   map[string]Function
	generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (c GoServiceConfig) Validate() bool {
	log.Printf("[INFO] VALIDATING GO SERVICE: %s", c.ServiceName)
	if !c.IsAlphaNumeric(c.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", c.ServiceName)
		return false
	}
	if c.RepoName == "" {
		log.Println("[ERROR] Repo name should not be empty")
		return false
	}
	for methodName, function := range c.Functions {
		if !c.IsAlphaNumeric(methodName) {
			log.Printf("[ERROR] method should be alphanumeric, but `%s` found", methodName)
			return false
		}
		if !function.Validate() {
			return false
		}
	}
	return true
}

// Scaffold scaffolding config
func (c GoServiceConfig) Scaffold() error {
	log.Printf("[INFO] SCAFFOLDING GO SERVICE: %s", c.ServiceName)
	for _, function := range c.Functions {
		data := function.ToExposed()
		packageSourcePath := function.FunctionPackage
		functionFileName := function.GetFunctionFileName()
		// write function
		functionSourcePath := filepath.Join(packageSourcePath, functionFileName)
		if !c.IsSourceExists(functionSourcePath) {
			log.Printf("[INFO] Create %s", functionFileName)
			err := c.WriteSource(functionSourcePath, "gosrvc.function.go", data)
			if err != nil {
				return err
			}
		} else {
			log.Printf("[INFO] %s exists", functionFileName)
		}
		// write dependencies
		for _, dependency := range function.FunctionDependencies {
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			if !c.IsSourceExists(dependencySourcePath) {
				log.Printf("[INFO] Create %s", dependency)
				err := c.WriteSource(dependencySourcePath, "dependency.go", data)
				if err != nil {
					return err
				}
			} else {
				log.Printf("[INFO] %s exists", dependency)
			}
		}
	}
	return nil
}

// Build building config
func (c GoServiceConfig) Build() error {
	log.Printf("[INFO] BUILDING GO SERVICE: %s", c.ServiceName)
	depPath := fmt.Sprintf("srvc-%s", c.ServiceName)
	repoName := c.RepoName
	mainFunctionName := "main"
	// create program
	err := c.CreateProgram(depPath, repoName, mainFunctionName)
	if err != nil {
		return err
	}
	// write common things
	err = c.WriteDeps(depPath, []string{"go.mod", "Makefile", ".gitignore"}, c)
	if err != nil {
		return err
	}
	// git init
	log.Printf("[INFO] Run git init")
	err = c.GitInitDep(depPath)
	if err != nil {
		return err
	}
	// GoFmt
	log.Printf("[INFO] Run gofmt")
	err = c.GoFmtDep(depPath)
	if err != nil {
		return err
	}
	return nil
}

// CreateProgram create main.go and others
func (c GoServiceConfig) CreateProgram(depPath, repoName, mainFunctionName string) error {
	// write functions and dependencies
	for _, function := range c.Functions {
		packageSourcePath := function.FunctionPackage
		packageDepPath := filepath.Join(depPath, function.FunctionPackage)
		functionFileName := function.GetFunctionFileName()
		// copy function
		log.Printf("[INFO] Create %s", functionFileName)
		functionSourcePath := filepath.Join(packageSourcePath, functionFileName)
		functionDepPath := filepath.Join(packageDepPath, functionFileName)
		err := c.CopySourceToDep(functionSourcePath, functionDepPath)
		if err != nil {
			return err
		}
		// copy dependencies
		for _, dependency := range function.FunctionDependencies {
			log.Printf("[INFO] Create %s", dependency)
			dependencySourcePath := filepath.Join(packageSourcePath, dependency)
			dependencyDepPath := filepath.Join(packageDepPath, dependency)
			err := c.CopySourceToDep(dependencySourcePath, dependencyDepPath)
			if err != nil {
				return err
			}
		}
	}
	// write main file
	mainFileName := fmt.Sprintf("%s.go", strings.ToLower(mainFunctionName))
	log.Printf("[INFO] Create %s", mainFileName)
	mainPath := filepath.Join(depPath, mainFileName)
	return c.WriteDep(mainPath, "gosrvc.main.go", c.toExposed(repoName, mainFunctionName))
}

// Set replace/add service's function
func (c *GoServiceConfig) Set(method string, function Function) {
	c.Functions[method] = function
}

func (c *GoServiceConfig) toExposed(repoName, mainFunctionName string) ExposedGoServiceConfig {
	exposedFunctions := make(map[string]ExposedFunction)
	var packages []string
	for methodName, function := range c.Functions {
		exposedFunction := function.ToExposed()
		exposedFunctions[methodName] = exposedFunction
		packages = append(packages, exposedFunction.FunctionPackage)
	}
	return ExposedGoServiceConfig{
		Packages:         packages,
		ServiceName:      c.ServiceName,
		RepoName:         repoName,
		MainFunctionName: mainFunctionName,
		Functions:        exposedFunctions,
	}
}

// NewGoService create new goservice
func NewGoService(g *generator.Generator, serviceName, repoName string, functions map[string]Function) GoServiceConfig {
	return GoServiceConfig{
		RepoName:    repoName,
		ServiceName: serviceName,
		Functions:   functions,
		IOHelper:    g.IOHelper,
	}
}

// NewEmptyGoService create new empty service
func NewEmptyGoService(g *generator.Generator, serviceName, repoName string) GoServiceConfig {
	return NewGoService(g, serviceName, repoName, make(map[string]Function))
}
