package msgbroker

import (
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"testing"
)

func TestNats(t *testing.T) {
	var broker CommonBroker
	broker, err := NewNats()
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	CommonBrokerTest(broker, t)
}

func TestNewCustomNats(t *testing.T) {
	_, err := NewCustomNats("invalid url")
	if err == nil {
		t.Error("Error expected")
	}
}

func TestNatsPublishInvalid(t *testing.T) {
	broker, err := NewNats()
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	// define package that cannot be marshaled into json
	pkg := servicedata.Package{
		ID: "some ID",
		Data: func(a, b int) int {
			return a + b
		},
	}
	// publish invalid package
	err = broker.Publish("invalidPublish", pkg)
	if err == nil {
		t.Error("Error expected")
	}
}

func TestNatsConsumeInvalid(t *testing.T) {
	eventName := "invalidConsume"
	broker, err := NewNats()
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	natsURL := GetNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	// consume
	stopped := make(chan bool, 1)
	broker.Consume(eventName,
		// consume success (should never happen)
		func(pkg servicedata.Package) {
			t.Errorf("Error expected")
			stopped <- true
		},
		// consume error (expected)
		func(err error) {
			if err == nil {
				t.Errorf("Error expected")
			}
			stopped <- true
		},
	)
	// publish
	err = nc.Publish(eventName, []byte("Hello world"))
	log.Print("publish")
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	<-stopped
}
