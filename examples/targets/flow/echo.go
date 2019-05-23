package main

import (
	"fmt"
	"github.com/state-alchemists/ayanami/service"
)

// Echo implemented by user
func Echo(text string) string {
	if text == "" {
		text = "nothing"
	}
	return fmt.Sprintf("You wrote %s", text)
}

// WrappedEcho implemented by user
func WrappedEcho(inputs service.Dictionary) service.Dictionary {
	outputs := make(service.Dictionary)
	// get text
	text := ExtractFormInterface(inputs["form"], "text")
	// prepare outputs
	outputs["content"] = Echo(text)
	outputs["code"] = 200
	return outputs
}
