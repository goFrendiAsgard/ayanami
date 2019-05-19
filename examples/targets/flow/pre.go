package main

import (
	"log"
)

// WrapPre call textDecorator.Pre
func WrapPre(inputs Dictionary) Dictionary {
	text := ExtractFormInterface(inputs["form"], "text")
	outputs, err := Call(
		"srvc", "textDecorator", "pre",
		[]string{"text"},
		[]string{"text"},
		Dictionary{
			"text": text,
		},
	)
	if err != nil {
		log.Printf("[ERROR] failed to call service: %s", err)
	}
	return outputs
}
