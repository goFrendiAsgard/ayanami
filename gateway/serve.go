package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type routeSorter []string

func (r routeSorter) Len() int {
	return len(r)
}

func (r routeSorter) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r routeSorter) Less(i, j int) bool {
	return len(r[i]) < len(r[j])
}

// Serve handle HTTP request
func Serve(broker msgbroker.CommonBroker, port int64, multipartFormLimit int64, routes []string) {
	sort.Sort(sort.Reverse(routeSorter(routes)))
	log.Printf("[INFO: Gateway] Routes `%#v`", routes)
	for _, route := range routes {
		handler := createRouteHandler(broker, multipartFormLimit, route)
		http.HandleFunc(route, handler)
	}
	log.Printf("Listening on %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func createRouteHandler(broker msgbroker.CommonBroker, multipartFormLimit int64, route string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// create ID
		ID, err := service.CreateID()
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		method := strings.ToLower(r.Method)
		// prepare channels
		codeChannel := make(chan int, 1)
		contentChannel := make(chan string, 1)
		consume(broker, ID, method, route, codeChannel, contentChannel)
		// publish
		err = publish(broker, ID, method, route, multipartFormLimit, r)
		if err != nil {
			responseError(ID, w, 500, err)
			return
		}
		// wait for response
		code := <-codeChannel
		content := <-contentChannel
		if code == 500 { // if there is a `500`, error, override the error message with this one
			content = "Internal Server Error"
		}
		response(ID, w, code, content)
	}
}

func consume(broker msgbroker.CommonBroker, ID, method, route string, codeChannel chan int, contentChannel chan string) {
	codeEventName := getResponseCodeEventName(ID, method, route)
	contentEventName := getResponseContentEventName(ID, method, route)
	// consume code
	log.Printf("[INFO: Gateway] Consume `%s`", codeEventName)
	broker.Consume(codeEventName,
		// success
		func(pkg servicedata.Package) {
			log.Printf("[INFO: Gateway] Getting message from `%s`: `%#v`", codeEventName, pkg.Data)
			val, err := strconv.Atoi(fmt.Sprintf("%v", pkg.Data))
			if err == nil {
				codeChannel <- val
			} else {
				log.Printf("[ERROR: Gateway] %s", err)
				codeChannel <- 500
				contentChannel <- "Internal Server Error"
			}
		},
		// error
		func(err error) {
			codeChannel <- 500
			contentChannel <- "Internal Server Error"
		},
	)
	// codeChannel <- 200
	// consume event
	log.Printf("[INFO: Gateway] Consume `%s`", contentEventName)
	broker.Consume(contentEventName,
		// success
		func(pkg servicedata.Package) {
			log.Printf("[INFO: Gateway] Getting message from `%s`: `%#v`", contentEventName, pkg.Data)
			contentChannel <- fmt.Sprintf("%s", pkg.Data)
		},
		// error
		func(err error) {
			codeChannel <- 500
			contentChannel <- "Internal Server Error"
		},
	)
}

func responseError(ID string, w http.ResponseWriter, code int, err error) {
	content := fmt.Sprintf("%s", err)
	response(ID, w, code, content)
}

func response(ID string, w http.ResponseWriter, code int, content string) {
	log.Printf("[INFO: Gateway] responding to %s: %d, %s", ID, code, content)
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", content)
}

func publish(broker msgbroker.CommonBroker, ID string, method string, route string, multipartFormLimit int64, r *http.Request) error {
	eventName := getRequestEventName(ID, method, route)
	// parse form & multipart form
	r.ParseForm()
	r.ParseMultipartForm(multipartFormLimit)
	data := make(map[string]interface{})
	data["header"] = r.Header
	data["contentLength"] = r.ContentLength
	data["host"] = r.Host
	data["form"] = r.Form
	data["postForm"] = r.PostForm
	data["multipartForm"] = r.MultipartForm
	data["method"] = r.Method
	data["requestURI"] = r.RequestURI
	data["remoteAddr"] = r.RemoteAddr
	// get json body
	decoder := json.NewDecoder(r.Body)
	JSONBody := make(map[string]interface{})
	err := decoder.Decode(&JSONBody)
	if err != nil {
		log.Printf("[INFO: Gateway] Processing `%s`, request.body is not a valid JSON", ID)
	}
	data["JSONBody"] = JSONBody
	// prepare pkg
	pkg := servicedata.Package{ID: ID, Data: data}
	log.Printf("[INFO: Gateway] publish `%s`: %#v", eventName, pkg)
	return broker.Publish(eventName, pkg)
}

func getResponseCodeEventName(ID, method, route string) string {
	return getEventName(ID, "response", method, route, "in", "code")
}

func getResponseContentEventName(ID, method, route string) string {
	return getEventName(ID, "response", method, route, "in", "content")
}

func getRequestEventName(ID, method, route string) string {
	return getEventName(ID, "request", method, route, "out", "req")
}

func getEventName(ID, trigger, method, route, direction, varName string) string {
	segment := RouteToSegments(route)
	if segment == "" {
		return fmt.Sprintf("%s.trig.%s.%s.%s.%s", ID, trigger, method, direction, varName)
	}
	return fmt.Sprintf("%s.trig.%s.%s.%s.%s.%s", ID, trigger, method, segment, direction, varName)
}
