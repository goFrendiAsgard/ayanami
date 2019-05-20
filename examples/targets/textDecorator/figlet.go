package main

import (
	"os/exec"
)

// Figlet implemented by user
func Figlet(text string) (string, error) {
	outByte, err := exec.Command("cowsay", text).Output()
	out := string(outByte)
	return out, err
}

// WrappedFiglet implemented by user
func WrappedFiglet(inputs SrvcDictionary) SrvcDictionary {
	outputs := make(SrvcDictionary)
	inputText := inputs["text"].(string)
	outputText, err := Figlet(inputText)
	outputs["text"] = outputText
	if err != nil {
		outputs["text"] = "figlet is not installed, here is your text: " + inputText
	}
	return outputs
}
