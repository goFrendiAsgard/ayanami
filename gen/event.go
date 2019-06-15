package gen

import (
	"fmt"
	"github.com/state-alchemists/ayanami/generator"
	"log"
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
	FunctionDependencies []string // should belong to the same package
	UseFunction          bool
	generator.StringHelper
}

// GetFunctionFileName get name of function file
func (e *Event) GetFunctionFileName() string {
	return fmt.Sprintf("%s.go", strings.ToLower(e.FunctionName))
}

// Validate validating an event
func (e *Event) Validate() bool {
	if e.VarName == "" {
		log.Println("[ERROR] Var name should not be empty")
		return false
	}
	if e.UseFunction && !e.IsMatch(e.FunctionName, "^[A-Z][a-zA-Z0-9]*$") && e.FunctionPackage == "" {
		log.Println("[ERROR] Function name should not be alphanumeric and function package should not be empty")
		return false
	}
	return true
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
		result["Function"] = fmt.Sprintf("%s.%s", e.FunctionPackage, e.FunctionName)
	}
	return result
}

// ToIndentedMap change event to map with value indented, ready for templating
func (e *Event) ToIndentedMap() map[string]string {
	result := e.ToMap()
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
