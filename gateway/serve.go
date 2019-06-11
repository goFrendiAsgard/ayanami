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

// Serve handle HTTP request
func Serve(broker msgbroker.CommonBroker, port int64, multipartFormLimit int64, routes []string) {
	sort.Sort(sort.Reverse(RouteSorter(routes)))
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
		method := strings.ToLower(r.Method)
		// create ID
		ID, err := service.CreateID()
		if err != nil {
			responseError(ID, broker, method, route, w, 500, err)
			return
		}
		// prepare channels
		codeChannel := make(chan int, 1)
		contentChannel := make(chan string, 1)
		consumeFromResponseTrigger(broker, ID, method, route, codeChannel, contentChannel)
		// publishToRequestTrigger
		err = publishToRequestTrigger(broker, ID, method, route, multipartFormLimit, r)
		if err != nil {
			responseError(ID, broker, method, route, w, 500, err)
			return
		}
		// wait for response
		code := <-codeChannel
		content := <-contentChannel
		if code == 500 { // if there is a `500`, error, override the error message with this one
			content = "Internal Server Error"
		}
		// response
		response(ID, broker, method, route, w, code, content)
	}
}

func consumeFromResponseTrigger(broker msgbroker.CommonBroker, ID, method, route string, codeChannel chan int, contentChannel chan string) {
	codeEventName := getResponseCodeEventName(ID, method, route)
	contentEventName := getResponseContentEventName(ID, method, route)
	// consumeFromResponseTrigger code
	log.Printf("[INFO: Gateway] Subscribe `%s`", codeEventName)
	broker.Subscribe(codeEventName,
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
	// consumeFromResponseTrigger event
	log.Printf("[INFO: Gateway] Subscribe `%s`", contentEventName)
	broker.Subscribe(contentEventName,
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

func responseError(ID string, broker msgbroker.CommonBroker, method, route string, w http.ResponseWriter, code int, err error) {
	content := fmt.Sprintf("%s", err)
	response(ID, broker, method, route, w, code, content)
}

func response(ID string, broker msgbroker.CommonBroker, method, route string, w http.ResponseWriter, code int, content string) {
	codeEventName := getResponseCodeEventName(ID, method, route)
	contentEventName := getResponseContentEventName(ID, method, route)
	err := broker.Unsubscribe(codeEventName)
	if err != nil {
		code = 500
		content = fmt.Sprintf("%s", err)
	}
	err = broker.Unsubscribe(contentEventName)
	if err != nil {
		code = 500
		content = fmt.Sprintf("%s", err)
	}
	log.Printf("[INFO: Gateway] responding to %s: %d, %s", ID, code, content)
	// TODO: User should be able to set their own content-types and other headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_, err = fmt.Fprintf(w, "%s", content)
	if err != nil {
		log.Printf("[ERROR: Gateway] responding to %s: %s", ID, err)
	}
}

func publishToRequestTrigger(broker msgbroker.CommonBroker, ID string, method string, route string, multipartFormLimit int64, r *http.Request) error {
	eventName := getRequestEventName(ID, method, route)
	data := getDataForPublish(ID, multipartFormLimit, r)
	return service.Publish("Gateway", "", broker, ID, eventName, data)
}

func getDataForPublish(ID string, multipartFormLimit int64, r *http.Request) map[string]interface{} {
	// parse form & multipart form
	err := r.ParseForm()
	if err != nil {
		log.Printf("[ERROR: Gateway] Processing `%s`: %s", ID, err)
	}
	err = r.ParseMultipartForm(multipartFormLimit)
	if err != nil {
		log.Printf("[ERROR: Gateway] Processing `%s`: %s", ID, err)
	}
	data := make(map[string]interface{})
	data["method"] = r.Method
	data["URL"] = r.URL
	data["proto"] = r.Proto
	data["protoMajor"] = r.ProtoMajor
	data["protoMinor"] = r.ProtoMinor
	data["header"] = r.Header
	data["contentLength"] = r.ContentLength
	data["transferEncoding"] = r.TransferEncoding
	data["host"] = r.Host
	data["form"] = r.Form
	data["postForm"] = r.PostForm
	data["multipartForm"] = r.MultipartForm
	data["trailer"] = r.Trailer
	data["remoteAddr"] = r.RemoteAddr
	data["requestURI"] = r.RequestURI
	data["cookies"] = r.Cookies()
	data["userAgent"] = r.UserAgent()
	// get json body
	decoder := json.NewDecoder(r.Body)
	JSONBody := make(map[string]interface{})
	err = decoder.Decode(&JSONBody)
	if err != nil {
		log.Printf("[INFO: Gateway] Processing `%s`, request.body is not a valid JSON", ID)
	}
	data["JSONBody"] = JSONBody
	return data
}

func getResponseCodeEventName(ID, method, route string) string {
	baseEventName := getBaseEventName(ID, "response", method, route, "in")
	return fmt.Sprintf("%s.%s", baseEventName, "code")
}

func getResponseContentEventName(ID, method, route string) string {
	baseEventName := getBaseEventName(ID, "response", method, route, "in")
	return fmt.Sprintf("%s.%s", baseEventName, "content")
}

func getRequestEventName(ID, method, route string) string {
	return getBaseEventName(ID, "request", method, route, "out")
}

func getBaseEventName(ID, trigger, method, route, direction string) string {
	segment := RouteToSegments(route)
	if segment == "" {
		return fmt.Sprintf("%s.trig.%s.%s.%s", ID, trigger, method, direction)
	}
	return fmt.Sprintf("%s.trig.%s.%s.%s.%s", ID, trigger, method, segment, direction)
}
