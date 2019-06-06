package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// CmdConfig configuration to generate Cmd
type CmdConfig struct {
	generator.Resource
}

// Validate validating config
func (config CmdConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config CmdConfig) Scaffold() error {
	return nil
}

// Build building config
func (config CmdConfig) Build() error {
	return nil
}
