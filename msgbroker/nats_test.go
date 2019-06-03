package msgbroker

import (
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"testing"
	"time"
)

func TestNats(t *testing.T) {
	var broker CommonBroker
	broker, err := NewNats(config.GetNatsURL())
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	CommonBrokerTest(broker, t)
}

func TestNatsInvalidURL(t *testing.T) {
	log.Println("Test invalid Nats URL")
	_, err := NewNats("invalid url")
	if err == nil {
		t.Error("Error expected")
	}
}

func TestNatsPublishInvalid(t *testing.T) {
	log.Println("Test invalid Nats publish event")
	broker, err := NewNats(config.GetNatsURL())
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
	log.Println("Test invalid Nats consume event")
	broker, err := NewNats(config.GetNatsURL())
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	natsURL := config.GetNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		t.Errorf("Get error: %s", err)
		return
	}
	// consume
	stopped := make(chan bool, 1)
	broker.Consume("invalidConsume",
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
	// if consume doesn't respond for too long, end it
	nc.Subscribe("invalidConsume", func(m *nats.Msg) {
		log.Printf("Get invalid consume package: %s", string(m.Data))
		time.Sleep(5 * time.Second)
		t.Errorf("Subscriber doesn't response for too long")
		stopped <- true
	})
	// publish
	err = nc.Publish("invalidConsume", []byte("Hello world"))
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
	<-stopped
}
