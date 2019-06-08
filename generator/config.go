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
func (cfgs Configs) Build() error {
	for _, config := range cfgs {
		err := config.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

// Scaffold scaffold configurations
func (cfgs Configs) Scaffold() error {
	for _, config := range cfgs {
		err := config.Scaffold()
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate validate configurations
func (cfgs Configs) Validate() bool {
	for _, config := range cfgs {
		if !config.Validate() {
			return false
		}
	}
	return true
}
