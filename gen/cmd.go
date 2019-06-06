package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"regexp"
)

// SingleCmd a definition of function
type SingleCmd struct {
	Inputs  []string
	Outputs []string
	Command string
}

// SingleCmdRep SingleCmd, ready to be parsed as template
type SingleCmdRep struct {
	Inputs  string
	Outputs string
	Command string
}

// CmdConfig configuration to generate Cmd
type CmdConfig struct {
	PackageName string
	ServiceName string
	Commands    map[string]SingleCmd
	*generator.Resource
}

// Validate validating config
func (config CmdConfig) Validate() bool {
	if config.PackageName == "" {
		log.Printf("[Invalid COmmand Service: %s] Package Name should not be empty", config.ServiceName)
		return false
	}
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumeric.Match([]byte(config.ServiceName)) {
		log.Printf("[Invalid Command Service: %s] Service name should be alphanumeric, but `%s` found", config.ServiceName, config.ServiceName)
		return false
	}
	for methodName, command := range config.Commands {
		if alphanumeric.Match([]byte(methodName)) {
			log.Printf("[Invalid Command Service: %s] method should be alphanumeric, but `%s` found", config.ServiceName, methodName)
			return false
		}
		if len(command.Inputs) == 0 {
			log.Printf("[Invalid Command Service: %s] command `%s` has no input", config.ServiceName, methodName)
			return false
		}
		if len(command.Outputs) == 0 {
			log.Printf("[Invalid Command Service: %s] command `%s` has no output", config.ServiceName, methodName)
			return false
		}
		if command.Command == "" {
			log.Printf("[Invalid Command Service: %s] command `%s` is empty", config.ServiceName, methodName)
			return false
		}
	}
	return true
}

// Scaffold scaffolding config
func (config CmdConfig) Scaffold() error {
	return nil
}

// Build building config
func (config CmdConfig) Build() error {
	// build template-ready commands
	tmplCommands := make(map[string]SingleCmdRep)
	for methodName, singleCmd := range config.Commands {
		singleCmdRep := SingleCmdRep{
			Inputs:  QuoteArray(singleCmd.Inputs, ", "),
			Outputs: QuoteArray(singleCmd.Outputs, ", "),
			Command: Quote(singleCmd.Command),
		}
		tmplCommands[methodName] = singleCmdRep
	}
	// write main.go
	mainPath := fmt.Sprintf("%s/main.go", config.ServiceName)
	err := config.WriteDep(mainPath, "cmd.go", tmplCommands)
	if err != nil {
		return err
	}
	// write go.mod
	goModPath := fmt.Sprintf("%s/go.mod", config.ServiceName)
	err = config.WriteDep(goModPath, "go.mod", config.PackageName)
	return err
}
