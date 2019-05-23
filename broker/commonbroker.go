package broker

import (
	"github.com/state-alchemists/ayanami/data"
)

// ConsumeFunc callback of broker's consumer
type ConsumeFunc func(pkg data.Package)

// CommonBroker interface of every broker
type CommonBroker interface {
	Consume(eventName string, callback ConsumeFunc)
	Publish(eventName string, pkg data.Package)
}
