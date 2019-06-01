package integrationtest

import (
	"github.com/state-alchemists/ayanami/gateway"
	"github.com/state-alchemists/ayanami/msgbroker"
	"log"
)

// MainGateway emulating gateway's main
func MainGateway() {
	routes := []string{
		"/",
	}
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	port := gateway.GetPort()
	multipartFormLimit := gateway.GetMultipartFormLimit()
	gateway.Serve(broker, port, multipartFormLimit, routes)
}
