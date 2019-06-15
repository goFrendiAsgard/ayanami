package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// BrokerNats represent monolith using nats broker
var BrokerNats = "nats"

// BrokerMemory represent monolith using memory broker
var BrokerMemory = "memory"

// ExposedGomonolithProc GomonolithProc for template
type ExposedGomonolithProc struct {
	Broker    string
	Functions []string
}

// GoMonolithProc procedureuration to generate GoMonolith
type GoMonolithProc struct {
	Broker   string
	DirName  string
	RepoName string
	generator.IOHelper
}

// Validate validating procedure
func (p GoMonolithProc) Validate(configs generator.Configs) bool {
	log.Printf("[INFO] VALIDATING GOMONOLITH: %s", p.DirName)
	if p.Broker != BrokerMemory && p.Broker != BrokerNats {
		log.Printf("[ERROR] Broker should be either `nats` or `memory`")
		return false
	}
	if p.DirName == "" {
		log.Printf("[ERROR] Dir name should not be empty")
		return false
	}
	if p.RepoName == "" {
		log.Printf("[ERROR] Repo name should not be empty")
		return false
	}
	return true
}

// Scaffold scaffolding procedure
func (p GoMonolithProc) Scaffold(configs generator.Configs) error {
	return nil
}

// Build building procedure
func (p GoMonolithProc) Build(configs generator.Configs) error {
	log.Printf("[INFO] BUILDING MONOLITH: %s", p.DirName)
	depPath := p.DirName
	var mainFunctionList []string
	for _, config := range configs {
		switch config.(type) {
		case CmdConfig:
			c := config.(CmdConfig)
			mainFunctionName := fmt.Sprintf("Service%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case GoServiceConfig:
			c := config.(GoServiceConfig)
			mainFunctionName := fmt.Sprintf("Service%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case GatewayConfig:
			c := config.(GatewayConfig)
			mainFunctionName := fmt.Sprintf("Gateway%s", strings.Title(c.ServiceName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		case FlowConfig:
			c := config.(FlowConfig)
			mainFunctionName := fmt.Sprintf("Flow%s", strings.Title(c.FlowName))
			mainFunctionList = append(mainFunctionList, mainFunctionName)
			err := c.CreateProgram(depPath, p.RepoName, mainFunctionName)
			if err != nil {
				return err
			}
		}
	}
	// write main.go
	data := ExposedGomonolithProc{Broker: p.Broker, Functions: mainFunctionList}
	mainPath := filepath.Join(depPath, "main.go")
	err := p.WriteDep(mainPath, "gomonolith.main.go", data)
	if err != nil {
		return err
	}
	// write common things
	for _, templateName := range []string{"go.mod", "Makefile", ".gitignore"} {
		log.Printf("[INFO] Create %s", templateName)
		goModPath := filepath.Join(depPath, templateName)
		err := p.WriteDep(goModPath, templateName, p)
		if err != nil {
			return err
		}
	}
	// git init
	gitPath := filepath.Join(depPath, ".git")
	if !p.IsDepExists(gitPath) {
		log.Printf("[INFO] Init git")
		shellCmd := exec.Command("git", "init")
		shellCmd.Dir = filepath.Join(p.GetDepPath(), depPath)
		err := shellCmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// NewGoMonolith make monolithic app
func NewGoMonolith(g *generator.Generator, broker, dirName, repoName string) GoMonolithProc {
	return GoMonolithProc{
		Broker:   broker,
		DirName:  dirName,
		RepoName: repoName,
		IOHelper: g.IOHelper,
	}
}
