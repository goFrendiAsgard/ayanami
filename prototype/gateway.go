package prototype

// GatewayConfig configuration to generate Gateway
type GatewayConfig struct {
	Routes []string
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
	return nil
}
