package generator

// CommonConfig common configuration interface
type CommonConfig interface {
	Build() error
	Scaffold() error
	Validate() bool
}

// Configs array of common configuration
type Configs []CommonConfig

// Build build configurations
func (configs Configs) Build() error {
	for _, config := range configs {
		err := config.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

// Scaffold scaffold configurations
func (configs Configs) Scaffold() error {
	for _, config := range configs {
		err := config.Scaffold()
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate validate configurations
func (configs Configs) Validate() bool {
	for _, config := range configs {
		if !config.Validate() {
			return false
		}
	}
	return true
}
