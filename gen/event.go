package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"strings"
)

// Event definition
type Event struct {
	InputEventName       string
	OutputEventName      string
	VarName              string
	Value                interface{}
	UseValue             bool
	FunctionName         string
	FunctionPackage      string
	FunctionDependencies []string
	UseFunction          bool
	generator.StringHelper
}

// ToMap change event to map
func (e *Event) ToMap() map[string]string {
	result := make(map[string]string)
	e.addQuotedToMapIfNotEmpty(result, "InputEvent", e.InputEventName)
	e.addQuotedToMapIfNotEmpty(result, "OutputEvent", e.OutputEventName)
	e.addQuotedToMapIfNotEmpty(result, "VarName", e.VarName)
	if e.UseValue {
		result["UseValue"] = "true"
		result["Value"] = fmt.Sprintf("%#v", e.Value)
	}
	if e.UseFunction {
		result["UseFunction"] = "true"
		packageParts := strings.Split(e.FunctionPackage, "/")
		namespace := packageParts[len(packageParts)-1]
		result["Function"] = fmt.Sprintf("%s.%s", namespace, e.FunctionName)
	}
	e.addValIndentationToMap(result)
	return result
}

func (e *Event) addQuotedToMapIfNotEmpty(m map[string]string, key, val string) {
	if val != "" {
		m[key] = e.Quote(val)
	}
}

func (e *Event) addValIndentationToMap(m map[string]string) {
	// get longest key
	maxKeyLength := 0
	for key := range m {
		keyLength := len(key)
		if keyLength > maxKeyLength {
			maxKeyLength = keyLength
		}
	}
	for key := range m {
		for i := len(key); i < maxKeyLength; i++ {
			m[key] = " " + m[key]
		}
	}
}

// NewInputEvent create new Event
func NewInputEvent(eventName, varName string) Event {
	return Event{InputEventName: eventName, VarName: varName}
}

// NewOutputEvent create new Event
func NewOutputEvent(eventName, varName string) Event {
	return Event{OutputEventName: eventName, VarName: varName}
}

// NewOutputEventVal create new Event with value
func NewOutputEventVal(eventName, varName string, value interface{}) Event {
	event := NewOutputEvent(eventName, varName)
	event.UseValue = true
	event.Value = value
	return event
}

// NewOutputEventFunc create new Event with function
func NewOutputEventFunc(eventName, varName, functionPackage, functionName string, functionDependencies []string) Event {
	event := NewOutputEvent(eventName, varName)
	event.UseFunction = true
	event.FunctionName = functionName
	event.FunctionPackage = functionPackage
	event.FunctionDependencies = functionDependencies
	return event
}
