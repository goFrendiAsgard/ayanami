package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
	"testing"
)

// CommonBrokerTest is general helper for testing commonBroker
func CommonBrokerTest(broker CommonBroker, t *testing.T) {
	sentPkg := servicedata.Package{ID: "001", Data: "Hello world"}
	// consume
	stopped := make(chan bool, 1)
	broker.Consume("test",
		// success
		func(pkg servicedata.Package) {
			if pkg.ID != "001" || pkg.Data != "Hello world" {
				t.Errorf("Expected `%#v`, get `%#v`", sentPkg, pkg)
			}
			stopped <- true
		},
		// error
		func(err error) {
			t.Errorf("Get error %s", err)
			stopped <- true
		},
	)
	// publish
	err := broker.Publish("test", sentPkg)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	<-stopped
}
