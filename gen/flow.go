package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// FlowEvent event definition
type FlowEvent struct {
	VarName              string
	InputEvent           string
	OutputEvent          string
	VarValue             string
	UseVarValue          bool
	FunctionName         string
	FunctionImportPath   string
	FunctionDependencies []string
	UseFunction          bool
}

// FlowConfig configuration to generate Flow
type FlowConfig struct {
	PackageName string
	FlowName    string
	Inputs      []string
	Outputs     []string
	Events      []FlowEvent
	*generator.IOHelper
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
