package main

import (
	"fmt"
)

// Pre implemented by user
func Pre(text string) string {
	return fmt.Sprintf("<pre>%s</pre>", text)
}

// WrappedPre implemented by user
func WrappedPre(inputs SrvcDictionary) SrvcDictionary {
	outputs := make(SrvcDictionary)
	inputText := inputs["text"].(string)
	outputText := Pre(inputText)
	outputs["text"] = outputText
	return outputs
}