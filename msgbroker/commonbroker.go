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
	Consume(eventName string, successCallback ConsumeSuccessFunc, errorCallback ConsumeErrorFunc)
	Publish(eventName string, pkg servicedata.Package) error
}
