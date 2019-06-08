package generator

// CommonProcedure procedure to build scaffold and validate
type CommonProcedure interface {
	Build(Configs) error
	Scaffold(Configs) error
	Validate(Configs) bool
}

// Procedures array of common procedures
type Procedures []CommonProcedure

// Build build configurations
func (procs Procedures) Build(configs Configs) error {
	for _, procedure := range procs {
		err := procedure.Build(configs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Scaffold scaffold configurations
func (procs Procedures) Scaffold(configs Configs) error {
	for _, procedure := range procs {
		err := procedure.Scaffold(configs)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate validate configurations
func (procs Procedures) Validate(configs Configs) bool {
	for _, procedure := range procs {
		if !procedure.Validate(configs) {
			return false
		}
	}
	return true
}
