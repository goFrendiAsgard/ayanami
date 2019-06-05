package prototype

// FlowConfig configuration to generate Flow
type FlowConfig struct {
}

// Validate validating config
func (config FlowConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config FlowConfig) Scaffold() error {
	return nil
}

// Build building config
func (config FlowConfig) Build() error {
	return nil
}
