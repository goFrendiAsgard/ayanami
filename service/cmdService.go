package service

import (
	"fmt"
	"os/exec"
)

// NewCmd create new cmd
func NewCmd(serviceName string, methodName string, inputs []string, outputs []string, command []string) CommonService {
	wrappedFunction := createCmdWrapper(command, inputs, outputs)
	return NewService(serviceName, methodName, inputs, outputs, wrappedFunction)
}

func createCmdWrapper(cmd []string, inputVarNames, outputVarNames []string) WrappedFunction {
	return func(inputs Dictionary) (Dictionary, error) {
		outputs := make(Dictionary)
		// get realCmd
		var realCmd []string
		for _, cmdPart := range cmd {
			isInput := false
			for _, inputVarName := range inputVarNames {
				if cmdPart == fmt.Sprintf("$%s", inputVarName) {
					isInput = true
					inputVal := fmt.Sprintf("%s", inputs.Get(inputVarName))
					realCmd = append(realCmd, inputVal)
				}
			}
			if !isInput {
				realCmd = append(realCmd, cmdPart)
			}
		}
		// run the command
		outByte, err := exec.Command(realCmd[0], realCmd[1:]...).Output()
		if err == nil {
			outputVal := string(outByte)
			for _, outputVarName := range outputVarNames {
				outputs[outputVarName] = outputVal
			}
		}
		return outputs, err
	}
}
