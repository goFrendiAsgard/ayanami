package main

import (
	"fmt"
)

// Echo implemented by user
func Echo(text string) string {
	if text == "" {
		text = "nothing"
	}
	return fmt.Sprintf("You wrote %s", text)
}

// WrappedEcho implemented by user
func WrappedEcho(inputs SrvcDictionary) SrvcDictionary {
	outputs := make(SrvcDictionary)
	// get text
	text := SrvcExtractFormInterface(inputs["form"], "text")
	// prepare outputs
	outputs["content"] = Echo(text)
	outputs["code"] = 200
	return outputs
}
