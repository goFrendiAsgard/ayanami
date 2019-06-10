package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
	"testing"
)

// CommonBrokerTest is general helper for testing commonBroker
func CommonBrokerTest(broker CommonBroker, t *testing.T) {
	pkg1 := servicedata.Package{ID: "001", Data: "Hello world"}
	pkg2 := servicedata.Package{ID: "002", Data: "Hi universe"}

	// consume 1
	stopped1 := make(chan bool, 1)
	broker.Subscribe("*.test.request.get.first.in",
		// success
		func(pkg servicedata.Package) {
			if pkg.ID != "001" || pkg.Data != "Hello world" {
				t.Errorf("Expected `%#v`, get `%#v`", pkg1, pkg)
			}
			stopped1 <- true
		},
		// error
		func(err error) {
			t.Errorf("Get error %s", err)
			stopped1 <- true
		},
	)
	// consume 2
	stopped2 := make(chan bool, 1)
	broker.Subscribe("ID.test.request.get.second.in",
		// success
		func(pkg servicedata.Package) {
			if pkg.ID != "002" || pkg.Data != "Hi universe" {
				t.Errorf("Expected `%#v`, get `%#v`", pkg2, pkg)
			}
			stopped2 <- true
		},
		// error
		func(err error) {
			t.Errorf("Get error %s", err)
			stopped2 <- true
		},
	)

	// publish 1
	err := broker.Publish("ID.test.request.get.first.in", pkg1)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	// publish 2
	err = broker.Publish("ID.test.request.get.second.in", pkg2)
	if err != nil {
		t.Errorf("Get error %s", err)
	}

	// wait
	<-stopped1
	<-stopped2
	err = broker.Unsubscribe("*.test.request.get.first.in")
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	err = broker.Unsubscribe("ID.test.request.get.second.in")
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	err = broker.Unsubscribe("oraono.event")
	if err == nil {
		t.Errorf("Error expected but get nil")
	}
}
