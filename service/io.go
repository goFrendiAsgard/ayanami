package service

import (
	"log"
	"strings"
)

// IO single IO configuration
type IO struct {
	VarName   string
	EventName string
}

// IOList list of IO
type IOList []IO

// GetUniqueVarNames get unique varNames from IO list
func (ioList IOList) GetUniqueVarNames() []string {
	var result []string
	for _, io := range ioList {
		result = AppendUniqueString(io.VarName, result)
	}
	return result
}

// GetTopLevelValues get dictionary with topLevelEventNames as key and it's values
func (ioList IOList) GetTopLevelValues(vars Dictionary) Dictionary {
	result := make(Dictionary)
	for _, io := range ioList {
		if !vars.Has(io.VarName) {
			continue
		}
		err := result.Set(io.EventName, vars.Get(io.VarName))
		if err != nil {
			log.Printf("[ERROR] Failed to set %s", io.EventName)
		}
	}
	return result
}

// GetTopLevelEventNames getting top level event names
func (ioList IOList) GetTopLevelEventNames() []string {
	var result []string
	for _, io := range ioList {
		eventParts := strings.Split(io.EventName, ".")
		var parentParts []string
		for index, part := range eventParts {
			parentParts = append(parentParts, part)
			if index > 2 && (part == "out" || part == "in" || part == "err") {
				break
			}
		}
		parentName := strings.Join(parentParts, ".")
		result = AppendUniqueString(parentName, result)
	}
	return result
}

// GetVarEventNames get eventNames from IO list with specified varName
func (ioList IOList) GetVarEventNames(varName string) []string {
	var result []string
	for _, io := range ioList {
		if io.VarName == varName {
			result = AppendUniqueString(io.EventName, result)
		}
	}
	return result
}

// GetUniqueEventNames get unique eventNames from IO list
func (ioList IOList) GetUniqueEventNames() []string {
	var result []string
	for _, io := range ioList {
		result = AppendUniqueString(io.EventName, result)
	}
	return result
}

// GetEventVarNames get varNames from IO list with specified eventName
func (ioList IOList) GetEventVarNames(eventName string) []string {
	var result []string
	for _, io := range ioList {
		if io.EventName == eventName {
			result = AppendUniqueString(io.VarName, result)
		}
	}
	return result
}
