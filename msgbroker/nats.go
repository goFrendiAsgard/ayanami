package msgbroker

import (
	"encoding/json"
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"os"
)

// Nats msgbroker
type Nats struct {
	nc *nats.Conn
}

// Consume nats.Consume
func (broker Nats) Consume(eventName string, callback ConsumeFunc) {
	broker.nc.Subscribe(eventName, func(m *nats.Msg) {
		var pkg servicedata.Package
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			log.Printf("[ERROR] Consuming %s from %s: %s", string(JSONByte), eventName, err)
			return
		}
		callback(pkg)
	})
}

// Publish nats.Publish
func (broker Nats) Publish(eventName string, pkg servicedata.Package) {
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] Publishing %#v to %s: %s", pkg, eventName, err)
		return
	}
	broker.nc.Publish(eventName, JSONByte)
}

// NewNats create new nats msgbroker
func NewNats() (CommonBroker, error) {
	var broker CommonBroker
	natsURL := getNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return broker, err
	}
	broker = Nats{nc: nc}
	return broker, err
}

func getNatsURL() string {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	return natsURL
}
