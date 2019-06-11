package msgbroker

import (
	"fmt"
	"github.com/state-alchemists/ayanami/servicedata"
	"regexp"
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
	// look for matched event
	for key, handler := range broker.handlers {
		// change eventName into regex
		eventPattern := fmt.Sprintf("^%s$", key)
		eventPattern = strings.Replace(eventPattern, ".", `\.`, -1)
		eventPattern = strings.Replace(eventPattern, "*", `[0-9a-zA-Z\*]+`, -1)
		eventPattern = strings.Replace(eventPattern, ">", ".*", -1)
		re, err := regexp.Compile(eventPattern)
		if err != nil {
			return err
		}
		if re.Match([]byte(eventName)) {
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
