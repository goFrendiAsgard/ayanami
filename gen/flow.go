package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// FlowConfig definition
type FlowConfig struct {
	PackageName  string
	FlowName     string
	Inputs       []string
	Outputs      []string
	InputEvents  []InputEvent
	OutputEvents []OutputEvent
	*generator.IOHelper
}

// Validate validating config
func (config *FlowConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config *FlowConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *FlowConfig) Build() error {
	return nil
}
