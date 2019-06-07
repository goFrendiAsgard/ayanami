package gen

import (
	"github.com/state-alchemists/ayanami/generator"
)

// FlowConfig definition
type FlowConfig struct {
	PackageName string
	FlowName    string
	Inputs      []string
	Outputs     []string
	Events      []Event
	*generator.IOHelper
}

func (config *FlowConfig) toMap() map[string]interface{} {
	result := make(map[string]interface{})
	// TODO see below
	// Packages []string <-- all event that use function and has dependencies
	// Events <-- array of map
	// Inputs string
	// Outputs string
	return result
}

// AddEvent add input to inputEvents
func (config *FlowConfig) AddEvent(event Event) {
	config.Events = append(config.Events, event)
}

// Validate validating config
func (config *FlowConfig) Validate() bool {
	return true
}

// Scaffold scaffolding config
func (config *FlowConfig) Scaffold() error {
	return nil
}

// Build building config
func (config *FlowConfig) Build() error {
	return nil
}

// NewFlow create new flow
func NewFlow(ioHelper *generator.IOHelper, packageName, flowName string, inputs, outputs []string, events []Event) FlowConfig {
	return FlowConfig{
		PackageName: packageName,
		FlowName:    flowName,
		Inputs:      inputs,
		Outputs:     outputs,
		Events:      events,
		IOHelper:    ioHelper,
	}
}

// NewEmptyFlow create new empty flow
func NewEmptyFlow(ioHelper *generator.IOHelper, packageName, flowName string, inputs, outputs []string) FlowConfig {
	return NewFlow(ioHelper, packageName, flowName, inputs, outputs, []Event{})
}
