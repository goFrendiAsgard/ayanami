package main

// SrvcEventFlow event flow
type SrvcEventFlow struct {
	InputEvent  string
	OutputEvent string
	Value       interface{} // if InputName == "", the Value will be published instead
}

// SrvcSingleFlowConfig single flow config
type SrvcSingleFlowConfig struct {
	FlowName string
	Input    []SrvcServiceIO
	Output   []SrvcServiceIO
	Flows    []SrvcEventFlow
}
