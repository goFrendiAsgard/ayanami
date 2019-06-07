package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// CmdConfig configuration to generate Cmd
type CmdConfig struct {
	ServiceName string
	PackageName string
	Commands    map[string]string
	*generator.IOHelper
	generator.StringHelper
}

// Set replace/add cmd's command
func (config *CmdConfig) Set(method, command string) {
	config.Commands[method] = command
}

// Validate validating config
func (config *CmdConfig) Validate() bool {
	log.Printf("[INFO] Validating %s", config.ServiceName)
	if config.IsAlphaNumeric(config.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", config.ServiceName)
		return false
	}
	if config.PackageName == "" {
		log.Println("[ERROR] Package name should not be empty")
		return false
	}
	for methodName, command := range config.Commands {
		if !config.IsAlphaNumeric(methodName) {
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
func (config *CmdConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *CmdConfig) Build() error {
	log.Printf("[INFO] Building %s", config.ServiceName)
	dirPath := fmt.Sprintf("srvc-%s", config.ServiceName)
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(dirPath, "main.go")
	err := config.WriteDep(mainPath, "cmd.main.go", config.QuoteMap(config.Commands))
	if err != nil {
		return err
	}
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(dirPath, "go.mod")
	err = config.WriteDep(goModPath, "go.mod", config.PackageName)
	return err
}

// NewCmd create new cmd
func NewCmd(ioHelper *generator.IOHelper, serviceName string, packageName string, commands map[string]string) CmdConfig {
	return CmdConfig{
		ServiceName: serviceName,
		PackageName: packageName,
		Commands:    commands,
		IOHelper:    ioHelper,
	}
}

// NewEmptyCmd create new empty cmd
func NewEmptyCmd(ioHelper *generator.IOHelper, serviceName string, packageName string) CmdConfig {
	return NewCmd(ioHelper, serviceName, packageName, make(map[string]string))
}
