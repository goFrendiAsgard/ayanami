package gateway

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func handle(broker msgbroker.CommonBroker, httpPort int, path string) {
	routes := []string{
		path,
	}
	port := int64(httpPort)
	multipartFormLimit := GetMultipartFormLimit()
	Serve(broker, port, multipartFormLimit, routes)
}

func handleTest(broker msgbroker.CommonBroker, port int, path string, t *testing.T) {
	broker.Consume(fmt.Sprintf("*.trig.request.get%s.out.req", RouteToSegments(path)),
		func(pkg servicedata.Package) {
			ID := pkg.ID
			// publish code
			codePkg := servicedata.Package{ID: ID, Data: 200}
			codeEvent := fmt.Sprintf("%s.trig.response.get%s.in.code", ID, RouteToSegments(path))
			broker.Publish(codeEvent, codePkg)
			// publish content
			contentPkg := servicedata.Package{ID: ID, Data: "hi"}
			contentEvent := fmt.Sprintf("%s.trig.response.get%s.in.content", ID, RouteToSegments(path))
			broker.Publish(contentEvent, contentPkg)
		},
		func(err error) {
			t.Errorf("Get error %s", err)
		},
	)
	go handle(broker, port, path)
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

func TestHandleWithMemory(t *testing.T) {
	broker, err := msgbroker.NewMemory()
	if err != nil {
		log.Fatal(err)
	}
	port := 8508
	handleTest(broker, port, "/memory", t)
}

func TestHandleWithNats(t *testing.T) {
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	port := 8507
	handleTest(broker, port, "/nats", t)
}
