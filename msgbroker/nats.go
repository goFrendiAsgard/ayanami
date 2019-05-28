package msgbroker

import (
	"encoding/json"
	nats "github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/servicedata"
	"os"
)

// Nats msgbroker
type Nats struct {
	Connection *nats.Conn
}

// Consume nats.Consume
func (broker *Nats) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.Connection.Subscribe(eventName, func(m *nats.Msg) {
		var pkg servicedata.Package
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			errorCallback(err)
			return
		}
		successCallback(pkg)
	})
}

// Publish nats.Publish
func (broker *Nats) Publish(eventName string, pkg servicedata.Package) error {
	// marshal package into JSON Byte
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		return err
	}
	return broker.Connection.Publish(eventName, JSONByte)
}

// NewNats create new nats msgbroker
func NewNats() (CommonBroker, error) {
	natsURL := getNatsURL()
	return NewCustomNats(natsURL)
}

// NewCustomNats create new nats msgbroker
func NewCustomNats(natsURL string) (CommonBroker, error) {
	var broker Nats
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return &broker, err
	}
	broker.Connection = nc
	return &broker, err
}

// getNatsURL get Nats URL from environment variable
func getNatsURL() string {
	// get natsURL from environment, or use defaultURL instead
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		natsURL = nats.DefaultURL
	}
	return natsURL
}