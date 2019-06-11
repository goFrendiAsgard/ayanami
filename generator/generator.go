package generator

import (
	"errors"
	"log"
)

// Generator used for scaffold and build
type Generator struct {
	configs    Configs
	procedures Procedures
	IOHelper   IOHelper
}

// AddConfig add single config to generator
func (g *Generator) AddConfig(config CommonConfig) {
	g.configs = append(g.configs, config)
}

// AddConfigs add single config to generator
func (g *Generator) AddConfigs(configs []CommonConfig) {
	for _, config := range configs {
		g.AddConfig(config)
	}
}

// AddProcedure add single procedure to generator
func (g *Generator) AddProcedure(procedure CommonProcedure) {
	g.procedures = append(g.procedures, procedure)
}

// AddProcedures add single procedure to generator
func (g *Generator) AddProcedures(procedures []CommonProcedure) {
	for _, procedure := range procedures {
		g.AddProcedure(procedure)
	}
}

// Validate validate all configs and procedures
func (g *Generator) Validate() bool {
	// validate all configs
	log.Println("VALIDATING CONFIGS")
	if !g.configs.Validate() {
		return false
	}
	// validate all procedures
	log.Println("VALIDATING PROCEDURES")
	if !g.procedures.Validate(g.configs) {
		return false
	}
	return true
}

// Build build from config & procedures
func (g *Generator) Build() error {
	// validate configs & procedures
	if !g.Validate() {
		return errors.New("invalid config/procedure")
	}
	// build configs
	log.Println("BUILDING CONFIGS")
	err := g.configs.Build()
	if err != nil {
		return err
	}
	// build procedures
	log.Println("BUILDING PROCEDURES")
	err = g.procedures.Build(g.configs)
	if err != nil {
		return err
	}
	return nil
}

// Scaffold scaffold from config & procedures
func (g *Generator) Scaffold() error {
	// validate configs & procedures
	if !g.Validate() {
		return errors.New("invalid config/procedure")
	}
	// scaffold configs
	log.Println("SCAFFOLDING CONFIGS")
	err := g.configs.Scaffold()
	if err != nil {
		return err
	}
	// scaffold procedures
	log.Println("SCAFFOLDING PROCEDURES")
	err = g.procedures.Scaffold(g.configs)
	if err != nil {
		return err
	}
	return nil
}

// NewGenerator create new generator
func NewGenerator() Generator {
	return Generator{}
}

// NewProjectGenerator create new generator
func NewProjectGenerator(projectPath string) (Generator, error) {
	generator := NewGenerator()
	ioHelper, err := NewIOHelperByProjectPath(projectPath)
	if err == nil {
		generator.IOHelper = ioHelper
	}
	return generator, err
}
