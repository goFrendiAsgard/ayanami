package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"reflect"
)

// Publish package
func Publish(serviceName, methodName string, broker msgbroker.CommonBroker, ID, eventName string, data interface{}) error {
	return publish(serviceName, methodName, broker, ID, eventName, data, 10)
}

func publish(serviceName, methodName string, broker msgbroker.CommonBroker, ID, eventName string, data interface{}, level int) error {
	if level <= 0 {
		return nil
	}
	pkg := servicedata.Package{ID: ID, Data: data}
	if methodName == "" {
		log.Printf("[INFO: %s] Publish `%s`: %#v", serviceName, eventName, data)
	} else {
		log.Printf("[INFO: %s.%s] Publish `%s`: %#v", serviceName, methodName, eventName, data)
	}
	err := broker.Publish(eventName, pkg)
	if err != nil {
		return err
	}
	reflectVal := reflect.ValueOf(data)
	reflectKind := reflectVal.Kind()
	// publish sub data
	if reflectKind == reflect.Map {
		for _, subVarName := range reflectVal.MapKeys() {
			subData := reflectVal.MapIndex(subVarName)
			subEventName := fmt.Sprintf("%s.%s", eventName, subVarName)
			err := publish(serviceName, methodName, broker, ID, subEventName, subData, level-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
