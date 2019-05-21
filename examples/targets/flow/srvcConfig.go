package main

// SrvcSingleConfig single configuration
type SrvcSingleConfig struct {
	Input    []SrvcServiceIO
	Output   []SrvcServiceIO
	Function SrvcWrappedFunction
}

// SrvcConfigs configuration
type SrvcConfigs = map[string]SrvcSingleConfig

// SrvcEventFlow event flow
type SrvcEventFlow struct {
	InputEvent  string
	OutputEvent string
}
