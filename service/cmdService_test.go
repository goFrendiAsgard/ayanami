package service

import (
	"reflect"
	"testing"
)

func TestNewCmdService(t *testing.T) {
	service := NewCmd("service", "method",
		[]string{"echo", "hello", "$name", "how are you", "${time}"},
	)

	// test inputs
	expectedInputs := IOList{
		IO{EventName: "srvc.service.method.in.name", VarName: "name"},
		IO{EventName: "srvc.service.method.in.time", VarName: "time"},
	}
	inputs := service.Input
	if !reflect.DeepEqual(inputs, expectedInputs) {
		t.Errorf("expected %#v, get %#v", expectedInputs, inputs)
	}

	// test outputs
	expectedOutputs := IOList{
		IO{EventName: "srvc.service.method.out.output", VarName: "output"},
	}
	Outputs := service.Output
	if !reflect.DeepEqual(Outputs, expectedOutputs) {
		t.Errorf("expected %#v, get %#v", expectedOutputs, Outputs)
	}

	// test errorEventName
	expectedErrorEventName := "srvc.service.method.err.message"
	if service.ErrorEventName != expectedErrorEventName {
		t.Errorf("expected %s, get %s", expectedErrorEventName, service.ErrorEventName)
	}

	// test wrappedFunction
	expectedFunctionOutput := make(Dictionary)
	expectedFunctionOutput["output"] = "hello world how are you today\n"
	functionInput := make(Dictionary)
	functionInput["name"] = "world"
	functionInput["time"] = "today"
	functionOutput, err := service.Function(functionInput)
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	if !reflect.DeepEqual(functionOutput, expectedFunctionOutput) {
		t.Errorf("expected %#v, get %#v", expectedFunctionOutput, functionOutput)
	}

}
