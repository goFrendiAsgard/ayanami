package msgbroker

import (
	"encoding/json"
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

// Unsubscribe unsubscribe from an event
func (broker Memory) Unsubscribe(eventName string) error {
	// check handler's existance
	broker.lock.RLock()
	_, exists := broker.handlers[eventName]
	broker.lock.RUnlock()
	// if handler is exist, remove it. Otherwise return error
	if exists {
		broker.lock.Lock()
		delete(broker.handlers, eventName)
		broker.lock.Unlock()
	} else {
		return fmt.Errorf("event `%s` doesn't exist", eventName)
	}
	return nil
}

// Publish publish to memory broker
func (broker Memory) Publish(eventName string, pkg servicedata.Package) error {
	data, err := json.Marshal(pkg.Data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, &pkg.Data); err != nil {
		return err
	}
	// get all keys of handler
	broker.lock.RLock()
	i := 0
	keys := make([]string, len(broker.handlers))
	for k := range broker.handlers {
		keys[i] = k
		i++
	}
	broker.lock.RUnlock()
	// look for matched event
	for _, key := range keys {
		// change eventName into regex
		eventPattern := fmt.Sprintf("^%s$", key)
		eventPattern = strings.Replace(eventPattern, ".", `\.`, -1)
		eventPattern = strings.Replace(eventPattern, "*", `[^\.]+`, -1)
		eventPattern = strings.Replace(eventPattern, ">", ".*", -1)
		re, err := regexp.Compile(eventPattern)
		if err != nil {
			return err
		}
		if re.Match([]byte(eventName)) {
			broker.lock.RLock()
			handler := broker.handlers[key]
			go handler(pkg)
			broker.lock.RUnlock()
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
