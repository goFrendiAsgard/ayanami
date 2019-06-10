package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
)

// ConsumeSuccessFunc called for every succeed consume
type ConsumeSuccessFunc func(pkg servicedata.Package)

// ConsumeErrorFunc called for every failed consume
type ConsumeErrorFunc func(err error)

// CommonBroker interface of every message broker
type CommonBroker interface {
	Subscribe(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc)
	Unsubscribe(eventName string) error
	Publish(eventName string, pkg servicedata.Package) error
}
