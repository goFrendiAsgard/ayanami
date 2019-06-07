package gen

// OutputEvent definition
type OutputEvent struct {
	EventName            string
	VarName              string
	Value                interface{}
	UseValue             bool
	FunctionName         string
	FunctionPackage      string
	FunctionDependencies []string
	UseFunction          bool
}

// NewOutputEvent create new OutputEvent
func NewOutputEvent(eventName, varName string) OutputEvent {
	return OutputEvent{EventName: eventName, VarName: varName}
}

// NewOutputEventVal create new OutputEvent with value
func NewOutputEventVal(eventName, varName string, value interface{}) OutputEvent {
	event := NewOutputEvent(eventName, varName)
	event.UseValue = true
	event.Value = value
	return event
}

// NewOutputEventFunc create new OutputEvent with function
func NewOutputEventFunc(eventName, varName, functionPackage, functionName string, functionDependencies []string) OutputEvent {
	event := NewOutputEvent(eventName, varName)
	event.UseFunction = true
	event.FunctionName = functionName
	event.FunctionPackage = functionPackage
	event.FunctionDependencies = functionDependencies
	return event
}
