# Ayanami

A FaaS-like framework for your own infrastructure.

The name is inspired from Evangelion-Unit-00's pilot: Ayanami Rei. The name `Rei` itself has the same pronunciation as in Heraclitus's philosophy, `Panta Rhei` (lit: everything flows). We believe that the developer should focus more on data flows and transformations rather than managing infrastructures.

# Why

* FaaS is good since it let you focus on the code instead of infrastructure
* Any FaaS providers are prone to vendor lock-in
* Having your own infrastructure (e.g: kubernetes) while developing/deploying in FaaS manner is a good solution
* At some point, developers need to run the entire infrastructure in their local machine. In this case, installing kubernetes/minikube could be overkill
* Generated instead of encapsulated

# Goal

Providing an environment with minimum dependencies in order to:

* Build & deploy FaaS
* Make kubernetes-ready artifacts
* Run the entire infrastructure locally without kubernetes

# Dependencies

* golang 1.2
* nats

# How

* Developer create functions. The functions can be written in any language, even binary
* Developer define flows (how the functions are connected to each others)
* Ayanami compose flows and functions into several microservices that can talk to each other using nats messaging.


# Terminologies

* Composition: The functionality definition of your program
    - Flow: Composition of functions, usually triggered by a trigger (e.g: when there is a HTTP request to `/order` end point, the system should execute several functions from different services and return a response).
    - Trigger: Event that trigger flows (e.g: Scheduler, HTTP request, etc)
    - Service: Collection of functions. Usually from the same domain
        - Function: The atomic part of your business logic. Functions from different services should be independent from each other.
* Template: The template we use to generate package
* Package: The final source code of your program, ready for deployment

# Convention

## Event Name

Event Name should comply one of these formats

```
<ID>.<trig|srvc|flow>.<serviceName>.<segments...>.<out|in>.<varName>
<ID>.<trig|srvc|flow>.<serviceName>.<segments...>.err.message
```

* `<ID>` is 32 characters of `UUID v4 with no hyphens`.
* `<trig|srvc|flow>` is service type, either `trig` (trigger), `srvc` (service), or `flow`.
* `<serviceName>` is either serviceName or flowname. Should only contains alphanumeric.
* `<segments...>` is description of the event. Should only contains alphanumeric or `.`, but should not started, ended, or has two consecutive `.`.
* `<out|in>` is either `out` or `in`. Typically services consume `in` event and omit `out` event.
* `<varName>` is variable name.

__Note:__ We strip `hyphens` from UUID because Nats documentation said it only accept alpha numeric and dots as event name.

# Gateway

Gateway provide two triggers:

* __Request trigger__ (`<ID>.trig.request.<http-verb>.<url-segments>.out.<request-var>`)
* __Response trigger__ (`<ID>.trig.response.<http-verb>.<url-segments>.in.<response-var>`)

Here are the valid values for each segments

* `<http-verb>` is lower case http verb. Either `post`, `get`, `put`, or `delete`
* `<url-segments>` is request URL with all `/` and spaces replaced into `.`
* `<request-var>` is a `map[string]interface{}`. It has several keys (please refer to golang's `net/http` documentation):
    - `header`
    - `contentLength`
    - `host`
    - `form`
    - `postForm`
    - `multipartForm`
    - `method`
    - `requestURI`
    - `remoteAddr`
    - `JSONBody`
* `responseVar` is either `code` (http status code) or `content`.

## Create Gateway

```go
package integrationtest

import (
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/gateway"
	"github.com/state-alchemists/ayanami/msgbroker"
	"log"
)

// MainGateway emulating gateway's main
func MainGateway() {
	routes := []string{ // define your routes
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
```

# Service
