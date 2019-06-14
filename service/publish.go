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

func getReflectAndRealData(data interface{}) (reflect.Value, reflect.Kind, interface{}) {
	reflectVal := reflect.ValueOf(data)
	reflectKind := reflectVal.Kind()
	if reflectKind == reflect.Ptr {
		reflectVal = reflectVal.Elem()
		reflectKind = reflectVal.Kind()
		if reflectVal.IsValid() {
			data = reflectVal.Interface()
		} else {
			data = nil
		}
	}
	return reflectVal, reflectKind, data
}

func getServicePlusMethodName(serviceName, methodName string) string {
	servicePlusMethod := serviceName
	if methodName != "" {
		servicePlusMethod = fmt.Sprintf("%s.%s", serviceName, methodName)
	}
	return servicePlusMethod
}

func publish(serviceName, methodName string, broker msgbroker.CommonBroker, ID, eventName string, data interface{}, level int) error {
	if level <= 0 {
		return nil
	}
	servicePlusMethod := getServicePlusMethodName(serviceName, methodName)
	reflectVal, reflectKind, data := getReflectAndRealData(data)
	// prepare package
	pkg := servicedata.Package{ID: ID, Data: data}
	log.Printf("[INFO: %s] Publish `%s`: %#v", servicePlusMethod, eventName, pkg)
	err := broker.Publish(eventName, pkg)
	if err != nil {
		return err
	}
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
	case reflect.Struct:
		for index := 0; index < reflectVal.NumField(); index++ {
			subData := reflectVal.Field(index).Interface()
			subVarName := reflectVal.Type().Field(index).Name
			subEventName := fmt.Sprintf("%s.%s", eventName, subVarName)
			err := publish(serviceName, methodName, broker, ID, subEventName, subData, level-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
