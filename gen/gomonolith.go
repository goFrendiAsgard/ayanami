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
func (p GoMonolithProc) Validate(configs generator.Configs) bool {
	return true
}

// Scaffold scaffolding procedure
func (p GoMonolithProc) Scaffold(configs generator.Configs) error {
	return nil
}

// Build building procedure
func (p GoMonolithProc) Build(configs generator.Configs) error {
	log.Printf("[INFO] Building %s", p.ServiceName)
	depPath := fmt.Sprintf(p.ServiceName)
	// write go.mod
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err := p.WriteDep(goModPath, "go.mod", p)
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
