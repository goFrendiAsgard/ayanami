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

// Consume consume from memory broker
func (broker Memory) Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc) {
	broker.lock.Lock()
	broker.handlers[eventName] = successCallback
	broker.lock.Unlock()
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
	broker = Memory{lock: &lock, handlers: handlers}
	return broker, nil
}
