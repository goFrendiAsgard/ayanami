package service

import (
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	service := NewService("service", "method",
		// inputs
		[]string{"a", "b"},
		// output
		[]string{"c"},
		// wrapped function
		func(inputs Dictionary) (Dictionary, error) {
			outputs := make(Dictionary)
			a := inputs.Get("a").(int)
			b := inputs.Get("b").(int)
			outputs["c"] = a + b
			return outputs, nil
		},
	)

	// test inputs
	expectedInputs := IOList{
		IO{EventName: "srvc.service.method.in.a", VarName: "a"},
		IO{EventName: "srvc.service.method.in.b", VarName: "b"},
	}
	inputs := service.Input
	if !reflect.DeepEqual(inputs, expectedInputs) {
		t.Errorf("expected %#v, get %#v", expectedInputs, inputs)
	}

	// test outputs
	expectedOutputs := IOList{
		IO{EventName: "srvc.service.method.out.c", VarName: "c"},
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
	expectedFunctionOutput["c"] = 9
	functionInput := make(Dictionary)
	functionInput["a"] = 4
	functionInput["b"] = 5
	functionOutput, err := service.Function(functionInput)
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	if !reflect.DeepEqual(functionOutput, expectedFunctionOutput) {
		t.Errorf("expected %#v, get %#v", expectedFunctionOutput, functionOutput)
	}

}
