package msgbroker

import (
	"testing"
)

func TestMemory(t *testing.T) {
	var broker CommonBroker
	broker, err := NewMemory()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	TestCommonBroker(broker, t)
}
