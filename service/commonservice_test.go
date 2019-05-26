package service

import (
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"testing"
)

func TestConsumeAndPublish(t *testing.T) {
	completed := make(chan bool)

	// define brokers
	broker, err := msgbroker.NewMemory()
	if err != nil {
		t.Errorf("Getting error: %s", err)
	}
	broker.Consume("*.srvc.common.out.c",
		func(pkg servicedata.Package) {
			c := pkg.Data.(int)
			if c != 7 {
				t.Errorf("expected 7, get %d", c)
			}
			completed <- true
		},
		func(err error) {
			t.Errorf("Getting error: %s", err)
			completed <- true
		},
	)
	broker.Consume("*.srvc.common.err",
		func(pkg servicedata.Package) {
			t.Errorf("Getting error: %s", pkg.Data)
			completed <- true
		},
		func(err error) {
			t.Errorf("Getting error: %s", err)
			completed <- true
		},
	)

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
	services.ConsumeAndPublish(broker)

	// publish a & b
	broker.Publish("ID.srvc.common.in.a", servicedata.Package{
		ID:   "ID",
		Data: 3,
	})
	broker.Publish("ID.srvc.common.in.b", servicedata.Package{
		ID:   "ID",
		Data: 4,
	})

	<-completed
}
