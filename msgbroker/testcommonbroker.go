package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
	"testing"
)

// TestCommonBroker is general helper for testing commonBroker
func TestCommonBroker(broker CommonBroker, t *testing.T) {
	sentPkg := servicedata.Package{ID: "001", Data: "Hello world"}
	// consume
	stopped := make(chan bool, 1)
	broker.Consume("test", func(pkg servicedata.Package) {
		if pkg.ID != "001" || pkg.Data != "Hello world" {
			t.Errorf("Expected `%#v`, get `%#v`", sentPkg, pkg)
		}
		stopped <- true
	})
	// publish
	broker.Publish("test", sentPkg)
	<-stopped
}
