package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// GoMonolithProc procedureuration to generate GoMonolith
type GoMonolithProc struct {
	ServiceName string
	RepoName    string
	generator.IOHelper
}

// Validate validating procedure
func (procedure GoMonolithProc) Validate(configs generator.Configs) bool {
	return true
}

// Scaffold scaffolding procedure
func (procedure GoMonolithProc) Scaffold(configs generator.Configs) error {
	return nil
}

// Build building procedure
func (procedure GoMonolithProc) Build(configs generator.Configs) error {
	log.Printf("[INFO] Building %s", procedure.ServiceName)
	depPath := fmt.Sprintf(procedure.ServiceName)
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err := procedure.WriteDep(goModPath, "go.mod", procedure)
	return err
}

// NewGoMonolith make monolithic app
func NewGoMonolith(g *generator.Generator, serviceName, repoName string) GoMonolithProc {
	return GoMonolithProc{
		ServiceName: serviceName,
		RepoName:    repoName,
		IOHelper:    g.IOHelper,
	}
}
