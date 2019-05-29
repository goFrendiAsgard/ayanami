package integrationtest

import (
	"fmt"
	"github.com/state-alchemists/ayanami/service"
)

// Pre implemented by user
func Pre(text string) string {
	return fmt.Sprintf("<pre>\n%s\n</pre>", text)
}

// WrappedPre implemented by user
func WrappedPre(inputs service.Dictionary) (service.Dictionary, error) {
	outputs := make(service.Dictionary)
	inputText := inputs["text"].(string)
	outputText := Pre(inputText)
	outputs["text"] = outputText
	return outputs, nil
}
