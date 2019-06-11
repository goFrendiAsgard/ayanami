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

func TestServe200WithMemory(t *testing.T) {
	broker, err := msgbroker.NewMemory()
	if err != nil {
		log.Fatal(err)
	}
	port := 8507
	serveTest200(broker, port, "/memory-200", t)
}

func TestServe500WithMemory(t *testing.T) {
	broker, err := msgbroker.NewMemory()
	if err != nil {
		log.Fatal(err)
	}
	port := 8508
	serveTest500(broker, port, "/memory-500", t)
}

func TestServeInvalidWithMemory(t *testing.T) {
	broker, err := msgbroker.NewMemory()
	if err != nil {
		log.Fatal(err)
	}
	port := 8509
	serveTestInvalidCode(broker, port, "/memory-invalid", t)
}

func TestServeWithNats(t *testing.T) {
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	port := 8511
	serveTest200(broker, port, "/nats", t)
}

func serveTest200(broker msgbroker.CommonBroker, port int, path string, t *testing.T) {
	broker.Subscribe(fmt.Sprintf("*.trig.request.get.%s.out", RouteToSegments(path)),
		func(pkg servicedata.Package) {
			// publish response
			ID := pkg.ID
			responsePkg := servicedata.Package{ID: ID, Data: map[string]interface{}{"code": 200, "content": "hi"}}
			eventName := fmt.Sprintf("%s.trig.response.get.%s.in", ID, RouteToSegments(path))
			err := broker.Publish(eventName, responsePkg)
			if err != nil {
				t.Errorf("Get error %s", err)
			}
		},
		func(err error) {
			t.Errorf("Get error %s", err)
		},
	)
	go serve(broker, port, path)
	time.Sleep(100 * time.Millisecond)
	// emulate request
	response, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, path))
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	// check body
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Get error %s", err)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	actualMessage := string(body)
	expectedMessage := "hi"
	if actualMessage != expectedMessage {
		t.Errorf("expectedMessage :%s, get :%s", expectedMessage, actualMessage)
	}
	// check code
	actualCode := response.StatusCode
	expectedCode := 200
	if actualCode != expectedCode {
		t.Errorf("expected: %d, get: %d", expectedCode, actualCode)
	}
}

func serveTest500(broker msgbroker.CommonBroker, port int, path string, t *testing.T) {
	broker.Subscribe(fmt.Sprintf("*.trig.request.get.%s.out", RouteToSegments(path)),
		func(pkg servicedata.Package) {
			// publish response
			ID := pkg.ID
			responsePkg := servicedata.Package{ID: ID, Data: map[string]interface{}{"code": 500, "content": "hi"}}
			eventName := fmt.Sprintf("%s.trig.response.get.%s.in", ID, RouteToSegments(path))
			err := broker.Publish(eventName, responsePkg)
			if err != nil {
				t.Errorf("Get error %s", err)
			}
		},
		func(err error) {
			t.Errorf("Get error %s", err)
		},
	)
	go serve(broker, port, path)
	time.Sleep(100 * time.Millisecond)
	// emulate request
	response, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, path))
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	// check body
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Get error %s", err)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	actualMessage := string(body)
	expectedMessage := "Internal Server Error"
	if actualMessage != expectedMessage {
		t.Errorf("expectedMessage :%s, get :%s", expectedMessage, actualMessage)
	}
	// check code
	actualCode := response.StatusCode
	expectedCode := 500
	if actualCode != expectedCode {
		t.Errorf("expected: %d, get: %d", expectedCode, actualCode)
	}
}

func serveTestInvalidCode(broker msgbroker.CommonBroker, port int, path string, t *testing.T) {
	broker.Subscribe(fmt.Sprintf("*.trig.request.get.%s.out", RouteToSegments(path)),
		func(pkg servicedata.Package) {
			// publishToRequestTrigger code
			ID := pkg.ID
			responsePkg := servicedata.Package{ID: ID, Data: map[string]interface{}{"code": "Not a valid code", "content": "hi"}}
			eventName := fmt.Sprintf("%s.trig.response.get.%s.in", ID, RouteToSegments(path))
			err := broker.Publish(eventName, responsePkg)
			if err != nil {
				t.Errorf("Get error %s", err)
			}
		},
		func(err error) {
			t.Errorf("Get error %s", err)
		},
	)
	go serve(broker, port, path)
	time.Sleep(100 * time.Millisecond)
	// emulate request
	response, err := http.Get(fmt.Sprintf("http://localhost:%d%s", port, path))
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	// check body
	defer func() {
		err := response.Body.Close()
		if err != nil {
			t.Errorf("Get error %s", err)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Get error %s", err)
	}
	actualMessage := string(body)
	expectedMessage := "Internal Server Error"
	if actualMessage != expectedMessage {
		t.Errorf("expectedMessage :%s, get :%s", expectedMessage, actualMessage)
	}
	// check code
	actualCode := response.StatusCode
	expectedCode := 500
	if actualCode != expectedCode {
		t.Errorf("expected: %d, get: %d", expectedCode, actualCode)
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
