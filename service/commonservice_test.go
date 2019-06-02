package service

import (
	"errors"
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"testing"
)

func TestConsumeAndPublish(t *testing.T) {
	// create broker
	broker, errorMessageCh, err := createCommonServiceBrokerTest("normal")
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// create service
	services := make(Services)
	services["test"] = CommonService{
		Input: IOList{
			IO{EventName: "srvc.common.in.a", VarName: "a"},
			IO{EventName: "srvc.common.in.b", VarName: "b"},
		},
		Output: IOList{
			IO{EventName: "srvc.common.out.c", VarName: "c"},
		},
		ErrorEventName: "srvc.common.err",
		Function: func(inputs Dictionary) (Dictionary, error) {
			outputs := make(Dictionary)
			a := inputs["a"].(int)
			b := inputs["b"].(int)
			outputs["c"] = a + b
			return outputs, nil
		},
	}
	services.ConsumeAndPublish(broker, "flow")
	// publish a & b
	broker.Publish("normal.srvc.common.in.a", servicedata.Package{
		ID:   "normal",
		Data: 3,
	})
	broker.Publish("normal.srvc.common.in.b", servicedata.Package{
		ID:   "normal",
		Data: 4,
	})
	errorMessage := <-errorMessageCh
	if errorMessage != "" {
		t.Errorf("Getting error: %s", errorMessage)
	}
}

func TestConsumeAndPublishFunctionError(t *testing.T) {
	// create broker
	broker, errorMessageCh, err := createCommonServiceBrokerTest("funcErr")
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// create service
	services := make(Services)
	services["test"] = CommonService{
		Input: IOList{
			IO{EventName: "srvc.common.in.a", VarName: "a"},
			IO{EventName: "srvc.common.in.b", VarName: "b"},
		},
		Output: IOList{
			IO{EventName: "srvc.common.out.c", VarName: "c"},
		},
		ErrorEventName: "srvc.common.err",
		Function: func(inputs Dictionary) (Dictionary, error) {
			outputs := make(Dictionary)
			return outputs, errors.New("ErrorThrown")
		},
	}
	services.ConsumeAndPublish(broker, "flow")
	// publish a & b
	broker.Publish("funcErr.srvc.common.in.a", servicedata.Package{
		ID:   "funcErr",
		Data: 3,
	})
	broker.Publish("funcErr.srvc.common.in.b", servicedata.Package{
		ID:   "funcErr",
		Data: 4,
	})
	errorMessage := <-errorMessageCh
	if errorMessage != "ErrorThrown" {
		t.Errorf("Expecting error message `ErrorThrown` getting `%s`", errorMessage)
	}
}

func createCommonServiceBrokerTest(ID string) (msgbroker.CommonBroker, chan string, error) {
	errorMessage := make(chan string)
	// define brokers
	broker, err := msgbroker.NewMemory()
	if err != nil {
		return broker, errorMessage, err
	}
	broker.Consume(fmt.Sprintf("%s.srvc.common.out.c", ID),
		func(pkg servicedata.Package) {
			c := pkg.Data.(int)
			if c != 7 {
				errorMessage <- fmt.Sprintf("expected 7, get %d", c)
			}
			errorMessage <- ""
		},
		func(err error) {
			errorMessage <- fmt.Sprintf("%s", err)
		},
	)
	broker.Consume(fmt.Sprintf("%s.srvc.common.err", ID),
		func(pkg servicedata.Package) {
			errorMessage <- fmt.Sprintf("%s", pkg.Data)
		},
		func(err error) {
			errorMessage <- fmt.Sprintf("%s", err)
		},
	)
	return broker, errorMessage, err
}
