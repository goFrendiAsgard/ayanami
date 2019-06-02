package gateway

import (
	"strings"
)

// RouteToSegments translate route into segments
func RouteToSegments(route string) string {
	route = strings.Replace(route, "/", ".", -1)
	route = strings.Replace(route, " ", ".", -1)
	if route == "." {
		route = ""
	}
	return route
}
