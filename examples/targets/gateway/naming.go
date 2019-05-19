package main

import (
	"fmt"
	"strings"
)

// DirectionOut output direction
const DirectionOut = "out"

// DirectionIn input direction
const DirectionIn = "in"

// ServiceTypeTrig trig service
const ServiceTypeTrig = "trig"

// ServiceTypeSrvc srvc service
const ServiceTypeSrvc = "srvc"

// ServiceTypeFlow flow service
const ServiceTypeFlow = "flow"

// GetEventName get event name
func GetEventName(ID string, serviceType string, srvcName string, fnName string, route string, direction string, varName string) string {
	return fmt.Sprintf("%s.%s.%s.%s.%s.%s.%s",
		ID,
		serviceType,
		srvcName,
		strings.ToLower(fnName),
		route,
		direction,
		varName,
	)
}
