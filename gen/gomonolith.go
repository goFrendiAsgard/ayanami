package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
	"strings"
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
	log.Printf("[INFO] BUILDING MONOLITH: %s", p.ServiceName)
	depPath := p.ServiceName
	mainFunctionList := []string{}
	for _, config := range configs {
		switch config.(type) {
		case CmdConfig:
			c := config.(CmdConfig)
			mainFunctionName := fmt.Sprintf("Service%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.ServiceName, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case GoServiceConfig:
			c := config.(GoServiceConfig)
			mainFunctionName := fmt.Sprintf("Service%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.ServiceName, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case GatewayConfig:
			c := config.(GatewayConfig)
			mainFunctionName := fmt.Sprintf("Gateway%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.ServiceName, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case FlowConfig:
			c := config.(FlowConfig)
			mainFunctionName := fmt.Sprintf("Flow%s", strings.Title(c.FlowName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.ServiceName, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		}
	}
	// TODO prepare to create main.go
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
