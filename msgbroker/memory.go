package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
)

// Memory broker for mocking
type Memory struct {
	handlers map[string]ConsumeSuccessFunc
}

// Consume consume from memory broker
func (broker *Memory) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.handlers[eventName] = successCallback
}

// Publish publish to memory broker
func (broker *Memory) Publish(eventName string, pkg servicedata.Package) error {
	go broker.handlers[eventName](pkg)
	return nil
}

// NewMemory create new memory brocker
func NewMemory() (CommonBroker, error) {
	handlers := make(map[string]ConsumeSuccessFunc)
	broker := Memory{handlers}
	return &broker, nil
}
