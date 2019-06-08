package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
	"path/filepath"
)

// ExposedGatewayConfig configuration to generate Gateway
type ExposedGatewayConfig struct {
	MainFunctionName string
	ServiceName      string
	RepoName         string
	Routes           []string
}

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	ServiceName string
	RepoName    string
	Routes      []string
	generator.IOHelper
	generator.StringHelper
}

// Validate validating config
func (c GatewayConfig) Validate() bool {
	log.Printf("[INFO] Validating %s", c.ServiceName)
	if !c.IsAlphaNumeric(c.ServiceName) {
		log.Printf("[ERROR] Service name should be alphanumeric, but `%s` found", c.ServiceName)
		return false
	}
	if c.RepoName == "" {
		log.Printf("[ERROR] Repo name should not be empty")
		return false
	}
	return true
}

// Scaffold scaffolding config
func (c GatewayConfig) Scaffold() error {
	log.Printf("[SKIP] Scaffolding %s", c.ServiceName)
	return nil
}

// Build building config
func (c GatewayConfig) Build() error {
	log.Printf("[INFO] Building %s", c.ServiceName)
	depPath := fmt.Sprintf("%s", c.ServiceName)
	// write main.go
	log.Println("[INFO] Create main.go")
	mainPath := filepath.Join(depPath, "main.go")
	err := c.WriteDep(mainPath, "gateway.main.go", c)
	if err != nil {
		return err
	}
	log.Println("[INFO] Create go.mod")
	goModPath := filepath.Join(depPath, "go.mod")
	err = c.WriteDep(goModPath, "go.mod", c)
	return err
}

// CreateProgram create main.go and others
func (c GatewayConfig) CreateProgram(depPath, serviceName, repoName, mainFunctionName string) error {
	// TODO use this
	return nil
}

// AddRoute add route to gateway
func (c *GatewayConfig) AddRoute(route string) {
	c.Routes = append(c.Routes, route)
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
