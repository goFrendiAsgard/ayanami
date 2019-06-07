package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// GoServiceConfig configuration to generate GoService
type GoServiceConfig struct {
	ServiceName string
	PackageName string
	Functions   map[string]Function
	*generator.IOHelper
}

// Validate validating config
func (config *GoServiceConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config *GoServiceConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *GoServiceConfig) Build() error {
	return nil
}
