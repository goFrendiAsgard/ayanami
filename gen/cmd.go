package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
	"strings"
)

// ExposedCmdConfig exposed ready cmdConfig
type ExposedCmdConfig struct {
	MainFunctionName string
	ServiceName      string
	RepoName         string
	Commands         map[string]string
}

// CmdConfig configuration to generate Cmd
type CmdConfig struct {
	ServiceName string
	RepoName    string
	Commands    map[string]string
	generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (c CmdConfig) Validate() bool {
	log.Printf("[INFO] VALIDATING CMD SERVICE: %s", c.ServiceName)
	if !c.IsAlphaNumeric(c.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", c.ServiceName)
		return false
	}
	if c.RepoName == "" {
		log.Println("[ERROR] Repo name should not be empty")
		return false
	}
	for methodName, command := range c.Commands {
		if !c.IsAlphaNumeric(methodName) {
			log.Printf("[ERROR] method should be alphanumeric, but `%s` found", methodName)
			return false
		}
		if command == "" {
			log.Printf("[ERROR] command `%s` is empty", methodName)
			return false
		}
	}
	return true
}

// Scaffold scaffolding config
func (c CmdConfig) Scaffold() error {
	return nil
}

// Build building config
func (c CmdConfig) Build() error {
	log.Printf("[INFO] BUILDING CMD SERVICE: %s", c.ServiceName)
	depPath := fmt.Sprintf("srvc-%s", c.ServiceName)
	serviceName := c.ServiceName
	repoName := c.RepoName
	mainFunctionName := "main"
	// create program
	err := c.CreateProgram(depPath, serviceName, repoName, mainFunctionName)
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
func (c CmdConfig) CreateProgram(depPath, serviceName, repoName, mainFunctionName string) error {
	// write main file
	mainFileName := fmt.Sprintf("%s.go", strings.ToLower(mainFunctionName))
	log.Printf("[INFO] Create %s", mainFileName)
	mainPath := filepath.Join(depPath, mainFileName)
	return c.WriteDep(mainPath, "cmd.main.go", c.toExposed(serviceName, repoName, mainFunctionName))
}

// Set replace/add cmd's command
func (c *CmdConfig) Set(method, command string) {
	c.Commands[method] = command
}

func (c *CmdConfig) toExposed(serviceName, repoName, mainFunctionName string) ExposedCmdConfig {
	return ExposedCmdConfig{
		ServiceName:      serviceName,
		RepoName:         repoName,
		MainFunctionName: mainFunctionName,
		Commands:         c.QuoteMap(c.Commands),
	}
}

// NewCmd create new cmd
func NewCmd(g *generator.Generator, serviceName string, repoName string, commands map[string]string) CmdConfig {
	return CmdConfig{
		ServiceName: serviceName,
		RepoName:    repoName,
		Commands:    commands,
		IOHelper:    g.IOHelper,
	}
}

// NewEmptyCmd create new empty cmd
func NewEmptyCmd(g *generator.Generator, serviceName string, repoName string) CmdConfig {
	return NewCmd(g, serviceName, repoName, make(map[string]string))
}
