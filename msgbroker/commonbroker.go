package msgbroker

import (
	"github.com/state-alchemists/ayanami/servicedata"
)

// ConsumeFunc callback of msgbroker's consumer
type ConsumeFunc func(pkg servicedata.Package)

// CommonBroker interface of every msgbroker
type CommonBroker interface {
	Consume(eventName string, callback ConsumeFunc)
	Publish(eventName string, pkg servicedata.Package)
}
