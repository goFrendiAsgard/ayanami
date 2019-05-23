package broker

import (
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/service"
	"log"
	"os"
)

// Nats broker
type Nats struct {
	nc *nats.Conn
}

// Consume nats.Consume
func (nats Nats) Consume(eventName string, callback ConsumeFunc) {
	nats.nc.Subscribe(eventName, func(m *nats.Msg) {
		var pkg service.Package
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			log.Printf("[ERROR] Consuming %s from %s: %s", string(JSONByte), eventName, err)
			return
		}
		consumeFunc(pkg)
	})
}

// Publish nats.Publish
func (nats Nats) Publish(eventName string, pkg service.Package) {
	JSONByte, _ := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] Publishing %#v to %s: %s", pkg, eventName, err)
		return
	}
	nats.nc.Publish(eventName, JSONByte)
}

// NewNats create new nats broker
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
