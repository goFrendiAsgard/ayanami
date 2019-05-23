package msgbroker

import (
	"testing"
)

func TestNats(t *testing.T) {
	var broker CommonBroker
	broker, err := NewNats()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	TestCommonBroker(broker, t)
}
