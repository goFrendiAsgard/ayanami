package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"reflect"
	"testing"
)

func TestFlowEvents(t *testing.T) {
	var expected, actual []string
	flowEvents := createFlowEventsTest()
	// getInputEvents
	expected = []string{"consume.a", "consume.b", "srvc.service.method.out.delta"}
	actual = flowEvents.GetInputEvents()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected `%#v`, get %#v`", expected, actual)
	}
	// getVarNamesByInputEvent
	expected = []string{"a"}
	actual = flowEvents.GetVarNamesByInputEvent("consume.a")
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected `%#v`, get %#v`", expected, actual)
	}
	// GetOutputEventByVarNames
	expected = []string{"srvc.service.method.in.alpha"}
	actual = flowEvents.GetOutputEventByVarNames("a")
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected `%#v`, get %#v`", expected, actual)
	}
}

func TestNewFlowService(t *testing.T) {
	broker := createFlowTestBroker(t)
	service := createFlowServiceTest(broker)
	// test inputs
	expectedInputs := IOList{
		IO{EventName: "flow.test.in.a", VarName: "a"},
		IO{EventName: "consume.a", VarName: "a"},
		IO{EventName: "flow.test.in.b", VarName: "b"},
		IO{EventName: "consume.b", VarName: "b"},
	}
	inputs := service.Input
	if !reflect.DeepEqual(inputs, expectedInputs) {
		t.Errorf("\nexpected: %#v\nget     : %#v", expectedInputs, inputs)
	}
	// test outputs
	expectedOutputs := IOList{
		IO{EventName: "flow.test.out.d", VarName: "d"},
		IO{EventName: "publish.d", VarName: "d"},
	}
	Outputs := service.Output
	if !reflect.DeepEqual(Outputs, expectedOutputs) {
		t.Errorf("\nexpected: %#v\nget     : %#v", expectedOutputs, Outputs)
	}
	// test errorEventName
	expectedErrorEventName := "flow.test.err.message"
	if service.ErrorEventName != expectedErrorEventName {
		t.Errorf("expected %s, get %s", expectedErrorEventName, service.ErrorEventName)
	}
	// test wrappedFunction
	expectedFunctionOutput := make(Dictionary)
	expectedFunctionOutput["d"] = 123
	functionInput := make(Dictionary)
	functionInput["a"] = 20
	functionInput["b"] = 3
	functionOutput, err := service.Function(functionInput)
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	if !reflect.DeepEqual(functionOutput, expectedFunctionOutput) {
		t.Errorf("expected %#v, get %#v", expectedFunctionOutput, functionOutput)
	}
}

func createFlowEventsTest() FlowEvents {
	return FlowEvents{
		FlowEvent{
			InputEvent:  "consume.a",
			VarName:     "a",
			OutputEvent: "srvc.service.method.in.alpha",
		},
		FlowEvent{
			InputEvent:  "consume.b",
			VarName:     "b",
			OutputEvent: "srvc.service.method.in.beta",
		},
		FlowEvent{
			VarName:     "c",
			Value:       100,
			OutputEvent: "srvc.service.method.in.gamma",
		},
		FlowEvent{
			InputEvent:  "srvc.service.method.out.delta",
			VarName:     "d",
			OutputEvent: "publish.d",
		},
	}
}

func createFlowServiceTest(broker msgbroker.CommonBroker) CommonService {
	// define flow
	service := NewFlow("flow", "test", broker,
		// inputs
		[]string{"a", "b"},
		// output
		[]string{"d"},
		// flows
		createFlowEventsTest(),
	)
	return service
}

func createFlowTestBroker(t *testing.T) msgbroker.CommonBroker {
	// prepare broker
	broker, err := msgbroker.NewMemory()
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	storage := NewDictionaryRW()
	errorCallback := func(err error) {
		storage.Set("a", 0)
		storage.Set("b", 0)
		storage.Set("c", 0)
		t.Errorf("Get error %s", err)
	}
	calculateAndPublish := func() {
		if storage.HasAll([]string{"a", "b", "c"}) {
			ID := storage.Get("ID").(string)
			a := storage.Get("a").(int)
			b := storage.Get("b").(int)
			c := storage.Get("c").(int)
			d := a + b + c
			pkg := servicedata.Package{ID: ID, Data: d}
			log.Printf("pkg: %#v\n", pkg)
			eventName := fmt.Sprintf("%s.srvc.service.method.out.delta", ID)
			broker.Publish(eventName, pkg)
		}
	}
	broker.Consume("*.srvc.service.method.in.alpha",
		func(pkg servicedata.Package) {
			storage.Set("ID", pkg.ID)
			storage.Set("a", pkg.Data.(int))
			calculateAndPublish()
		},
		errorCallback,
	)
	broker.Consume("*.srvc.service.method.in.beta",
		func(pkg servicedata.Package) {
			storage.Set("ID", pkg.ID)
			storage.Set("b", pkg.Data.(int))
			calculateAndPublish()
		},
		errorCallback,
	)
	broker.Consume("*.srvc.service.method.in.gamma",
		func(pkg servicedata.Package) {
			storage.Set("ID", pkg.ID)
			storage.Set("c", pkg.Data.(int))
			calculateAndPublish()
		},
		errorCallback,
	)
	broker.Consume("*.publish.d",
		func(pkg servicedata.Package) {
		},
		errorCallback,
	)
	return broker
}
