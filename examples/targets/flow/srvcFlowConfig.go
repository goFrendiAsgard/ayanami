package main

// SrvcEventFlow event flow
type SrvcEventFlow struct {
	InputEvent  string
	OutputEvent string
	VarName     string      // read from inputEvent, put into var, publish into outputEvent
	Value       interface{} // if InputEvent == "", the Value will be published instead
}

// SrvcSingleFlowConfig single flow config
type SrvcSingleFlowConfig struct {
	FlowName string
	Input    []SrvcServiceIO
	Output   []SrvcServiceIO
	Flows    []SrvcEventFlow
}
