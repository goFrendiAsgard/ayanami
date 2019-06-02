package service

import (
	"fmt"
)

// NewService create new singleConfig for service
func NewService(serviceName string, methodName string, inputs, outputs []string, wrappedFunction WrappedFunction) CommonService {
	// get serviceInputs
	var serviceInputs []IO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, inputName)
		serviceInputs = append(serviceInputs, IO{VarName: inputName, EventName: eventName})
	}
	// get serviceOutputs
	var serviceOutputs []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("srvc.%s.%s.out.%s", serviceName, methodName, outputName)
		serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
	}
	// get errorEventName
	errorEventName := fmt.Sprintf("srvc.%s.%s.err.message", serviceName, methodName)
	// return config
	return CommonService{
		Input:          serviceInputs,
		Output:         serviceOutputs,
		ErrorEventName: errorEventName,
		Function:       wrappedFunction,
	}
}
