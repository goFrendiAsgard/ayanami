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
	errorMessageCh := make(chan string, 1)
	broker, err := createCommonServiceBrokerTest("normal", errorMessageCh)
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// create service
	services := Services{
		CommonService{
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
	errorMessageCh := make(chan string, 1)
	broker, err := createCommonServiceBrokerTest("funcErr", errorMessageCh)
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	// create service
	services := Services{
		CommonService{
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

func createCommonServiceBrokerTest(ID string, errorMessageCh chan string) (msgbroker.CommonBroker, error) {
	// define brokers
	broker, err := msgbroker.NewMemory()
	if err != nil {
		return broker, err
	}
	// consume event
	broker.Consume(fmt.Sprintf("%s.srvc.common.out.c", ID),
		func(pkg servicedata.Package) {
			c := pkg.Data.(int)
			if c != 7 {
				errorMessageCh <- fmt.Sprintf("expected 7, get %d", c)
			}
			errorMessageCh <- ""
		},
		func(err error) {
			errorMessageCh <- fmt.Sprintf("%s", err)
		},
	)
	// consume error event
	broker.Consume(fmt.Sprintf("%s.srvc.common.err", ID),
		func(pkg servicedata.Package) {
			errorMessageCh <- fmt.Sprintf("%s", pkg.Data)
		},
		func(err error) {
			errorMessageCh <- fmt.Sprintf("%s", err)
		},
	)
	return broker, err
}
