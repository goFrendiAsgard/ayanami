package generator

import (
	"errors"
)

// Generator used for scaffold and build
type Generator struct {
	configs    Configs
	procedures Procedures
}

// AddConfig add single config to generator
func (generator *Generator) AddConfig(config CommonConfig) {
	generator.configs = append(generator.configs, config)
}

// AddConfigs add single config to generator
func (generator *Generator) AddConfigs(configs []CommonConfig) {
	for _, config := range configs {
		generator.AddConfig(config)
	}
}

// AddProcedure add single procedure to generator
func (generator *Generator) AddProcedure(procedure CommonProcedure) {
	generator.procedures = append(generator.procedures, procedure)
}

// AddProcedures add single procedure to generator
func (generator *Generator) AddProcedures(procedures []CommonProcedure) {
	for _, procedure := range procedures {
		generator.AddProcedure(procedure)
	}
}

// Validate validate all configs and procedures
func (generator *Generator) Validate() bool {
	// validate all configs
	if !generator.configs.Validate() {
		return false
	}
	// validate all procedures
	if !generator.procedures.Validate(generator.configs) {
		return false
	}
	return true
}

// Build build from config & procedures
func (generator *Generator) Build() error {
	// validate configs & procedures
	if !generator.Validate() {
		return errors.New("Invalid config/procedure")
	}
	// build configs
	err := generator.configs.Build()
	if err != nil {
		return err
	}
	// build procedures
	err = generator.procedures.Build(generator.configs)
	if err != nil {
		return err
	}
	return nil
}

// Scaffold scaffold from config & procedures
func (generator *Generator) Scaffold() error {
	// validate configs & procedures
	if !generator.Validate() {
		return errors.New("Invalid config/procedure")
	}
	// scaffold configs
	err := generator.configs.Scaffold()
	if err != nil {
		return err
	}
	// scaffold procedures
	err = generator.procedures.Scaffold(generator.configs)
	if err != nil {
		return err
	}
	return nil
}
