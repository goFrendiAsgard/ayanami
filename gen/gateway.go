package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	Routes []string
	generator.Resource
}

// Validate validating config
func (config GatewayConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config GatewayConfig) Scaffold() error {
	return nil
}

// Build building config
func (config GatewayConfig) Build() error {
	return config.Resource.WriteDeployable("gateway.go", "gateway.go", config)
}
