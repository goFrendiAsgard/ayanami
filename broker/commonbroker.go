package broker

import(
	"github.com/state-alchemists/ayanami/Package"
)

type CommonBroker interface{
	Consume(eventName string, pkg Package)
	Publish()
}
