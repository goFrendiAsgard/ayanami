package service

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// NewCmd create new cmd
func NewCmd(serviceName string, methodName string, inputs []string, outputs []string, command []string) CommonService {
	wrappedFunction := createCmdWrapper(serviceName, methodName, command, inputs, outputs)
	return NewService(serviceName, methodName, inputs, outputs, wrappedFunction)
}

func createCmdWrapper(serviceName, methodName string, cmd []string, inputVarNames, outputVarNames []string) WrappedFunction {
	return func(inputs Dictionary) (Dictionary, error) {
		outputs := make(Dictionary)
		// preprocess cmd
		for cmdIndex := range cmd {
			for _, varName := range inputVarNames {
				pattern1 := fmt.Sprintf("$%s", varName)
				pattern2 := fmt.Sprintf("${%s}", varName)
				varValue := fmt.Sprintf("%s", inputs.Get(varName))
				// if varValue doesn't started and ended with double quote, add double quote to it. Otherwise, let it be
				if cmd[cmdIndex] != pattern1 && cmd[cmdIndex] != pattern2 {
					valParts := strings.Split(varValue, "")
					if valParts[0] != "\"" || valParts[len(valParts)-1] != "\"" {
						varValue = strings.Replace(varValue, "\"", "\\\"", -1)
						varValue = fmt.Sprintf("\"%s\"", varValue)
					}
				}
				cmd[cmdIndex] = strings.Replace(cmd[cmdIndex], pattern1, varValue, -1)
				cmd[cmdIndex] = strings.Replace(cmd[cmdIndex], pattern2, varValue, -1)
			}
		}
		// run the command
		log.Printf("[INFO: %s.%s] Executing `%#v`", serviceName, methodName, cmd)
		outByte, err := exec.Command(cmd[0], cmd[1:]...).Output()
		if err != nil {
			log.Printf("[ERROR: %s.%s] Error while executing `%#v`: %s", serviceName, methodName, cmd, err)
			return outputs, err
		}
		// assemble outputs
		outputVal := string(outByte)
		for _, outputVarName := range outputVarNames {
			outputs.Set(outputVarName, outputVal)
		}
		return outputs, err
	}
}
