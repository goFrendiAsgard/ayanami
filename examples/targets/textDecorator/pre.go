package main

import (
	"fmt"
	"github.com/state-alchemists/ayanami/service"
)

// Pre implemented by user
func Pre(text string) string {
	return fmt.Sprintf("<pre>%s</pre>", text)
}

// WrappedPre implemented by user
func WrappedPre(inputs service.Dictionary) service.Dictionary {
	outputs := make(service.Dictionary)
	inputText := inputs["text"].(string)
	outputText := Pre(inputText)
	outputs["text"] = outputText
	return outputs
}
