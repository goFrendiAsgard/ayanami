package msgbroker

import (
	"testing"
)

func TestMemory(t *testing.T) {
	var broker CommonBroker
	broker, err := NewMemory()
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	CommonBrokerTest(broker, "ID.test", "ID.test", t)
	CommonBrokerTest(broker, "*.testWildcard", "ID.testWildcard", t)
}
