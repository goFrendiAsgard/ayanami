package gen

import (
	"github.com/state-alchemists/ayanami/generator"
	"log"
)

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	ServiceName string
	PackageName string
	Routes      []string
	*generator.IOHelper
}

// Validate validating config
func (config *GatewayConfig) Validate() bool {
	if config.PackageName == "" {
		log.Printf("[Invalid Gateway] Package Name should not be empty")
		return false
	}
	return true
}

// Scaffold scaffolding config
func (config *GatewayConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *GatewayConfig) Build() error {
	err := config.WriteDep("gateway/main.go", "gateway.main.go", config)
	if err != nil {
		return err
	}
	err = config.WriteDep("gateway/go.mod", "go.mod", config.PackageName)
	return err
}

// NewGateway create new gateway
func NewGateway(ioHelper *generator.IOHelper, serviceName string, packageName string, routes []string) GatewayConfig {
	return GatewayConfig{
		ServiceName: serviceName,
		PackageName: packageName,
		Routes:      routes,
		IOHelper:    ioHelper,
	}
}
