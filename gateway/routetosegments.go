package gateway

import (
	"regexp"
)

// RouteToSegments translate route into segments
func RouteToSegments(route string) string {
	// replace all forbidden characters into "."
	re := regexp.MustCompile(`[/ ]+`)
	route = re.ReplaceAllLiteralString(route, ".")
	// normalize all consecutive "."
	re = regexp.MustCompile(`\.+`)
	route = re.ReplaceAllLiteralString(route, ".")
	// remove "." at the begining and end of string
	re = regexp.MustCompile(`^\.`)
	route = re.ReplaceAllLiteralString(route, "")
	re = regexp.MustCompile(`\.$`)
	route = re.ReplaceAllLiteralString(route, "")
	return route
}
