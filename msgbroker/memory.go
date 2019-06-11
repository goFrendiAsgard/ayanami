package msgbroker

import (
	"fmt"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"strings"
	"sync"
)

// Memory broker for mocking
type Memory struct {
	lock     *sync.RWMutex
	handlers map[string]ConsumeSuccessFunc
}

// Subscribe consume from memory broker
func (broker Memory) Subscribe(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.lock.Lock()
	broker.handlers[eventName] = successCallback
	broker.lock.Unlock()
}

// Unsubscribe unsubscribe to an event
func (broker Memory) Unsubscribe(eventName string) error {
	broker.lock.Lock()
	if _, exists := broker.handlers[eventName]; exists {
		delete(broker.handlers, eventName)
	} else {
		return fmt.Errorf("event `%s` doesn't exist", eventName)
	}
	broker.lock.Unlock()
	return nil
}

// Publish publish to memory broker
func (broker Memory) Publish(eventName string, pkg servicedata.Package) error {
	broker.lock.RLock()
	defer broker.lock.RUnlock()
	// log.Printf("[MEMORY PUBLISH]\n  Event  : %s\n  Content: %#v", eventName, pkg)
	if handler, exists := broker.handlers[eventName]; exists {
		log.Printf("[MEMORY CONSUME]\n  Event  : %s\n  Content: %#v", eventName, pkg)
		go handler(pkg)
		return nil
	}
	eventParts := strings.Split(eventName, ".")
	wildCardEventName := fmt.Sprintf("*.%s", strings.Join(eventParts[1:], "."))
	handler, exists := broker.handlers[wildCardEventName]
	if exists {
		// log.Printf("[MEMORY CONSUME]\n  Event  : %s\n  Content: %#v", wildCardEventName, pkg)
		go handler(pkg)
	}
	return nil
}

// NewMemory create new memory brocker
func NewMemory() (CommonBroker, error) {
	var broker CommonBroker
	handlers := make(map[string]ConsumeSuccessFunc)
	lock := sync.RWMutex{}
	broker = Memory{lock: &lock, handlers: handlers}
	return broker, nil
}
