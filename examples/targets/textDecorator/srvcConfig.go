package main

// SrvcSingleConfig single configuration
type SrvcSingleConfig struct {
	Input    []SrvcServiceIO
	Output   []SrvcServiceIO
	Function SrvcWrappedFunction
}

// SrvcConfigs configuration
type SrvcConfigs = map[string]SrvcSingleConfig
