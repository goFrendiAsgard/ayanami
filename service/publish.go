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
	servicePlusMethod := serviceName
	if methodName != "" {
		servicePlusMethod = fmt.Sprintf("%s.%s", serviceName, methodName)
	}
	log.Printf("[INFO: %s] Publish `%s`: %#v", servicePlusMethod, eventName, pkg)
	err := broker.Publish(eventName, pkg)
	if err != nil {
		return err
	}
	reflectVal := reflect.ValueOf(data)
	reflectKind := reflectVal.Kind()
	// publish sub data recursively
	switch reflectKind {
	case reflect.Map:
		iter := reflectVal.MapRange()
		for iter.Next() {
			subVarName := iter.Key().Interface()
			subData := iter.Value().Interface()
			subEventName := fmt.Sprintf("%s.%s", eventName, subVarName)
			err := publish(serviceName, methodName, broker, ID, subEventName, subData, level-1)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		for index := 0; index < reflectVal.Len(); index++ {
			subData := reflectVal.Index(index).Interface()
			subEventName := fmt.Sprintf("%s.%d", eventName, index)
			err := publish(serviceName, methodName, broker, ID, subEventName, subData, level-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
