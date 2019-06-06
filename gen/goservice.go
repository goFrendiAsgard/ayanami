package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// Function a definition of function
type Function struct {
	Inputs       []string
	Outputs      []string
	Namespace    string
	Name         string
	Dependencies []string
	Template     string
}

// GoServiceConfig configuration to generate GoService
type GoServiceConfig struct {
	PackageName string
	*generator.Resource
}

// Validate validating config
func (config GoServiceConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config GoServiceConfig) Scaffold() error {
	return nil
}

// Build building config
func (config GoServiceConfig) Build() error {
	return nil
}
