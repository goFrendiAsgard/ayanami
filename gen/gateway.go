package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	ServiceName string
	RepoName    string
	Routes      []string
	generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (config GatewayConfig) Validate() bool {
	log.Printf("[INFO] Validating %s", config.ServiceName)
	if !config.IsAlphaNumeric(config.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", config.ServiceName)
		return false
	}
	if config.RepoName == "" {
		log.Printf("[ERROR] Repo name should not be empty")
		return false
	}
	return true
}

// Scaffold scaffolding config
func (config GatewayConfig) Scaffold() error {
	log.Printf("[SKIP] Scaffolding %s", config.ServiceName)
	return nil
}

// Build building config
func (config GatewayConfig) Build() error {
	log.Printf("[INFO] Building %s", config.ServiceName)
	depPath := fmt.Sprintf("%s", config.ServiceName)
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(depPath, "main.go")
	err := config.WriteDep(mainPath, "gateway.main.go", config)
	if err != nil {
		return err
	}
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err = config.WriteDep(goModPath, "go.mod", config)
	return err
}

// AddRoute add route to gateway
func (config *GatewayConfig) AddRoute(route string) {
	config.Routes = append(config.Routes, route)
}

// NewGateway create new gateway
func NewGateway(g *generator.Generator, serviceName string, repoName string, routes []string) GatewayConfig {
	return GatewayConfig{
		ServiceName: serviceName,
		RepoName:    repoName,
		Routes:      routes,
		IOHelper:    g.IOHelper,
	}
}

// NewEmptyGateway create new empty gateway
func NewEmptyGateway(g *generator.Generator, serviceName string, repoName string) GatewayConfig {
	return NewGateway(g, serviceName, repoName, []string{})
}
