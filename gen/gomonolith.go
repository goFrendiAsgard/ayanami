package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// GoMonolithProcedure procedureuration to generate GoMonolith
type GoMonolithProcedure struct {
	*generator.Resource
}

// Validate validating procedure
func (procedure GoMonolithProcedure) Validate(config generator.CommonConfig) bool {
	return true
}

// Scaffold scaffolding procedure
func (procedure GoMonolithProcedure) Scaffold(config generator.CommonConfig) error {
	return nil
}

// Build building procedure
func (procedure GoMonolithProcedure) Build(config generator.CommonConfig) error {
	return nil
}
