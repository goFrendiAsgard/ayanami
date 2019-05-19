package main

import (
	"log"
)

// WrapCowsay call textDecorator.Cowsay
func WrapCowsay(inputs Dictionary) Dictionary {
	text := ExtractFormInterface(inputs["form"], "text")
	outputs, err := Call(
		"srvc", "textDecorator", "cowsay",
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
