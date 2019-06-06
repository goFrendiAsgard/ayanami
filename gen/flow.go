package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// FlowConfig configuration to generate Flow
type FlowConfig struct {
	PackageName string
	*generator.Resource
}

// Validate validating config
func (config FlowConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config FlowConfig) Scaffold() error {
	return nil
}

// Build building config
func (config FlowConfig) Build() error {
	return nil
}
