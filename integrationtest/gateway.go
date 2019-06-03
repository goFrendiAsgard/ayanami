package integrationtest

import (
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/gateway"
	"github.com/state-alchemists/ayanami/msgbroker"
	"log"
)

// MainGateway emulating gateway's main
func MainGateway() {
	routes := []string{
		"/",
	}
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	port := config.GetGatewayPort()
	multipartFormLimit := config.GetGatewayMultipartFormLimit()
	gateway.Serve(broker, port, multipartFormLimit, routes)
}
