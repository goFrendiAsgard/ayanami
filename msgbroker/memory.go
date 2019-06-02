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
	handlers map[string]ConsumeSuccessFunc
	lock     *sync.RWMutex
}

// Consume consume from memory broker
func (broker Memory) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.lock.Lock()
	defer broker.lock.Unlock()
	broker.handlers[eventName] = successCallback
}

// Publish publish to memory broker
func (broker Memory) Publish(eventName string, pkg servicedata.Package) error {
	broker.lock.RLock()
	defer broker.lock.RUnlock()
	log.Printf("[MEMORY PUBLISH]\n  Event  : %s\n  Content: %#v", eventName, pkg)
	if handler, exists := broker.handlers[eventName]; exists {
		log.Printf("[MEMORY CONSUME]\n  Event  : %s\n  Content: %#v", eventName, pkg)
		go handler(pkg)
	} else {
		eventParts := strings.Split(eventName, ".")
		wildCardEventName := fmt.Sprintf("*.%s", strings.Join(eventParts[1:], "."))
		handler, exists := broker.handlers[wildCardEventName]
		if exists {
			log.Printf("[MEMORY CONSUME]\n  Event  : %s\n  Content: %#v", wildCardEventName, pkg)
			go handler(pkg)
		}
	}
	return nil
}

// NewMemory create new memory brocker
func NewMemory() (CommonBroker, error) {
	var broker CommonBroker
	handlers := make(map[string]ConsumeSuccessFunc)
	lock := sync.RWMutex{}
	broker = Memory{handlers, &lock}
	return broker, nil
}
