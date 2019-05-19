package main

import (
	"os/exec"
)

// Cowsay implemented by user
func Cowsay(text string) (string, error) {
	outByte, err := exec.Command("cowsay", text).Output()
	out := string(outByte)
	return out, err
}

// WrappedCowsay implemented by user
func WrappedCowsay(inputs Dictionary) Dictionary {
	outputs := make(Dictionary)
	inputText := inputs["text"].(string)
	outputText, err := Cowsay(inputText)
	outputs["text"] = outputText
	if err != nil {
		outputs["text"] = "cowsay is not installed, here is your text: " + inputText
	}
	return outputs
}
