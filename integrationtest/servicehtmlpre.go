package integrationtest

import (
	"fmt"
	"github.com/state-alchemists/ayanami/service"
)

// Pre implemented by user
func Pre(text string) string {
	return fmt.Sprintf("<pre>%s</pre>", text)
}

// WrappedPre implemented by user
func WrappedPre(inputs service.Dictionary) (service.Dictionary, error) {
	outputs := make(service.Dictionary)
	inputText := inputs["input"].(string)
	outputText := Pre(inputText)
	outputs["output"] = outputText
	return outputs, nil
}
