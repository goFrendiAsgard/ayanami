package msgbroker

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"sync"
)

// Nats msgbroker
type Nats struct {
	Connection    *nats.Conn
	lock          *sync.RWMutex
	subscriptions map[string]*nats.Subscription
}

// Subscribe nats.Subscribe
func (broker Nats) Subscribe(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	subscription, err := broker.Connection.Subscribe(eventName, func(m *nats.Msg) {
		pkg := servicedata.Package{}
		JSONByte := m.Data
		err := json.Unmarshal(JSONByte, &pkg)
		if err != nil {
			log.Printf("[NATS ERROR]\n  Event  : %s\n  Error: %s", eventName, err)
			errorCallback(err)
			return
		}
		successCallback(pkg)
	})
	if err != nil {
		log.Printf("[NATS ERROR]\n  Event  : %s\n  Error: %#v", eventName, err)
		errorCallback(err)
	}
	broker.lock.Lock()
	broker.subscriptions[eventName] = subscription
	broker.lock.Unlock()
}

// Unsubscribe unsubscribe to an event
func (broker Nats) Unsubscribe(eventName string) error {
	broker.lock.Lock()
	if subscription, exists := broker.subscriptions[eventName]; exists {
		err := subscription.Drain()
		if err != nil {
			return err
		}
		err = subscription.Unsubscribe()
		if err != nil {
			return err
		}
		delete(broker.subscriptions, eventName)
	} else {
		return fmt.Errorf("event `%s` doesn't exist", eventName)
	}
	broker.lock.Unlock()
	return nil
}

// Publish nats.Publish
func (broker Nats) Publish(eventName string, pkg servicedata.Package) error {
	// marshal package into JSON Byte
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		return err
	}
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
	broker.lock = &sync.RWMutex{}
	broker.subscriptions = make(map[string]*nats.Subscription)
	return broker, err
}
