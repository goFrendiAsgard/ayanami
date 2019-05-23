package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
)

// Memory broker for mocking
type Memory struct {
	handlers map[string]ConsumeFunc
}

// Consume consume from memory broker
func (broker *Memory) Consume(eventName string, callback ConsumeFunc) {
	broker.handlers[eventName] = callback
}

// Publish publish to memory broker
func (broker *Memory) Publish(eventName string, pkg servicedata.Package) {
	go broker.handlers[eventName](pkg)
}

// NewMemory create new memory brocker
func NewMemory() (CommonBroker, error) {
	handlers := make(map[string]ConsumeFunc)
	broker := Memory{handlers}
	return &broker, nil
}
