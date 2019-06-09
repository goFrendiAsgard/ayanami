package msgbroker

import (
	"encoding/json"
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
)

// Nats msgbroker
type Nats struct {
	Connection *nats.Conn
}

// Consume nats.Consume
func (broker Nats) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	_, err := broker.Connection.Subscribe(eventName, func(m *nats.Msg) {
		pkg := servicedata.Package{}
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			log.Printf("[NATS ERROR]\n  Event  : %s\n  Error: %s", eventName, err)
			errorCallback(err)
			return
		}
		// log.Printf("[NATS CONSUME]\n  Event  : %s\n  Content: %#v", eventName, pkg)
		successCallback(pkg)
	})
	if err != nil {
		log.Printf("[NATS ERROR]\n  Event  : %s\n  Error: %#v", eventName, err)
		errorCallback(err)
	}
}

// Publish nats.Publish
func (broker Nats) Publish(eventName string, pkg servicedata.Package) error {
	// marshal package into JSON Byte
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		return err
	}
	// log.Printf("[NATS PUBLISH]\n  Event  : %s\n  Content: %#v", eventName, pkg)
	return broker.Connection.Publish(eventName, JSONByte)
}

// NewNats create new nats msgbroker
func NewNats(natsURL string) (CommonBroker, error) {
	var broker Nats
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return broker, err
	}
	broker.Connection = nc
	return broker, err
}
