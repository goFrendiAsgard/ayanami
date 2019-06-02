package msgbroker

import (
	"fmt"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"strings"
)

// Memory broker for mocking
type Memory struct {
	handlers map[string]ConsumeSuccessFunc
}

// Consume consume from memory broker
func (broker Memory) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.handlers[eventName] = successCallback
}

// Publish publish to memory broker
func (broker Memory) Publish(eventName string, pkg servicedata.Package) error {
	log.Printf("[MEMORY PUBLISH]\nEvent  : %s\nContent: %#v", eventName, pkg)
	if handler, exists := broker.handlers[eventName]; exists {
		log.Printf("[MEMORY CONSUME]\nEvent  : %s\nContent: %#v", eventName, pkg)
		go handler(pkg)
	} else {
		eventParts := strings.Split(eventName, ".")
		wildCardEventName := fmt.Sprintf("*.%s", strings.Join(eventParts[1:], "."))
		handler, exists := broker.handlers[wildCardEventName]
		if exists {
			log.Printf("[MEMORY CONSUME]\nEvent  : %s\nContent: %#v", wildCardEventName, pkg)
			go handler(pkg)
		}
	}
	return nil
}

// NewMemory create new memory brocker
func NewMemory() (CommonBroker, error) {
	var broker CommonBroker
	handlers := make(map[string]ConsumeSuccessFunc)
	broker = Memory{handlers}
	return broker, nil
}
