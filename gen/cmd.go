package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
)

// Cmd a definition of function
type Cmd struct {
	Inputs  []string
	Outputs []string
	Command string
}

// CmdConfig configuration to generate Cmd
type CmdConfig struct {
	ServiceName string
	PackageName string
	Commands    map[string]Cmd
	*generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (config CmdConfig) Validate() bool {
	if config.PackageName == "" {
		log.Println("[ERROR] Package Name should not be empty")
		return false
	}
	if config.IsAlphaNumeric(config.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", config.ServiceName)
		return false
	}
	for methodName, command := range config.Commands {
		if config.IsAlphaNumeric(methodName) {
			log.Printf("[ERROR] method should be alphanumeric, but `%s` found", methodName)
			return false
		}
		if len(command.Inputs) == 0 {
			log.Printf("[ERROR] command `%s` has no input", methodName)
			return false
		}
		if len(command.Outputs) == 0 {
			log.Printf("[ERROR] command `%s` has no output", methodName)
			return false
		}
		if command.Command == "" {
			log.Printf("[ERROR] command `%s` is empty", methodName)
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
	tmplCommands := make(map[string]map[string]string)
	for methodName, cmd := range config.Commands {
		tmplCommands[methodName] = map[string]string{
			"Inputs":  config.QuoteArrayAndJoin(cmd.Inputs, ", "),
			"Outputs": config.QuoteArrayAndJoin(cmd.Outputs, ", "),
			"Command": config.Quote(cmd.Command),
		}
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
