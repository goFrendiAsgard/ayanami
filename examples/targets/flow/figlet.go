package main

import (
	"log"
)

// WrapFiglet call textDecorator.Figlet
func WrapFiglet(inputs Dictionary) Dictionary {
	text := ExtractFormInterface(inputs["form"], "text")
	outputs, err := Call(
		"srvc", "textDecorator", "figlet",
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
