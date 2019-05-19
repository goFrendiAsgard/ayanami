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
func WrappedEcho(inputs Dictionary) Dictionary {
	outputs := make(Dictionary)
	// get text
	form := inputs["form"].(map[string]interface{})
	text := ""
	if rawTexts, exists := form["text"].([]interface{}); exists {
		rawText := rawTexts[0]
		text = rawText.(string)
	}
	// prepare outputs
	outputs["content"] = Echo(text)
	outputs["code"] = 200
	return outputs
}
