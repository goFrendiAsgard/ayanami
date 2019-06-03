package gateway

import (
	"fmt"
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestServeWithMemory(t *testing.T) {
	broker, err := msgbroker.NewMemory()
	if err != nil {
		log.Fatal(err)
	}
	port := 8508
	serveTest(broker, port, "/memory", t)
}

func TestServeWithNats(t *testing.T) {
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	port := 8507
	serveTest(broker, port, "/nats", t)
}

func serveTest(broker msgbroker.CommonBroker, port int, path string, t *testing.T) {
	broker.Consume(fmt.Sprintf("*.trig.request.get.%s.out.req", RouteToSegments(path)),
		func(pkg servicedata.Package) {
			ID := pkg.ID
			// publish code
			codePkg := servicedata.Package{ID: ID, Data: 200}
			codeEvent := fmt.Sprintf("%s.trig.response.get.%s.in.code", ID, RouteToSegments(path))
			broker.Publish(codeEvent, codePkg)
			// publish content
			contentPkg := servicedata.Package{ID: ID, Data: "hi"}
			contentEvent := fmt.Sprintf("%s.trig.response.get.%s.in.content", ID, RouteToSegments(path))
			broker.Publish(contentEvent, contentPkg)
		},
		func(err error) {
			t.Errorf("Get error %s", err)
		},
	)
	go serve(broker, port, path)
	time.Sleep(1 * time.Second)
	// emulate request
	response, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, path))
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	actual := string(body)
	expected := "hi"
	if actual != expected {
		t.Errorf("expected :\n%s, get :\n%s", expected, actual)
	}
}

func serve(broker msgbroker.CommonBroker, httpPort int, path string) {
	routes := []string{
		path,
	}
	port := int64(httpPort)
	multipartFormLimit := config.GetGatewayMultipartFormLimit()
	Serve(broker, port, multipartFormLimit, routes)
}
