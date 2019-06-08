package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// ExposedGoServiceConfig exposed ready flowConfig
type ExposedGoServiceConfig struct {
	Packages    []string
	ServiceName string
	RepoName    string
	Functions   map[string]ExposedFunction
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
	log.Printf("[INFO] Validating %s", c.ServiceName)
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
	log.Printf("[INFO] Scaffolding %s", c.ServiceName)
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
	log.Printf("[INFO] Building %s", c.ServiceName)
	depPath := fmt.Sprintf("srvc-%s", c.ServiceName)
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
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(depPath, "main.go")
	err := c.WriteDep(mainPath, "gosrvc.main.go", c.toExposed())
	if err != nil {
		return err
	}
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err = c.WriteDep(goModPath, "go.mod", c)
	return err
}

// CreateProgram create main.go and others
func (c GoServiceConfig) CreateProgram(depPath, serviceName, repoName, mainFunction string) {
	// TODO use this
}

// Set replace/add service's function
func (c *GoServiceConfig) Set(method string, function Function) {
	c.Functions[method] = function
}

func (c *GoServiceConfig) toExposed() ExposedGoServiceConfig {
	exposedFunctions := make(map[string]ExposedFunction)
	packages := []string{}
	for methodName, function := range c.Functions {
		exposedFunction := function.ToExposed()
		exposedFunctions[methodName] = exposedFunction
		packages = append(packages, exposedFunction.FunctionPackage)
	}
	return ExposedGoServiceConfig{
		Packages:    packages,
		RepoName:    c.RepoName,
		ServiceName: c.ServiceName,
		Functions:   exposedFunctions,
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
