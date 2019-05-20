package main

// SrvcServiceIO single IO configuration
type SrvcServiceIO struct {
	VarName   string
	EventName string
}

// SrvcGetUniqueVarNames get unique varNames from SrvcServiceIO list
func SrvcGetUniqueVarNames(ioList []SrvcServiceIO) []string {
	var result []string
	for _, io := range ioList {
		if !isStringInArray(io.VarName, result) {
			result = append(result, io.VarName)
		}
	}
	return result
}

// SrvcGetVarEventNames get eventNames from SrvcServiceIO list with specified varName
func SrvcGetVarEventNames(ioList []SrvcServiceIO, varName string) []string {
	var result []string
	for _, io := range ioList {
		if io.VarName == varName {
			result = append(result, io.EventName)
		}
	}
	return result
}

// SrvcGetUniqueEventNames get unique eventNames from SrvcServiceIO list
func SrvcGetUniqueEventNames(ioList []SrvcServiceIO) []string {
	var result []string
	for _, io := range ioList {
		if !isStringInArray(io.EventName, result) {
			result = append(result, io.EventName)
		}
	}
	return result
}

// SrvcGetEventVarNames get varNames from SrvcServiceIO list with specified eventName
func SrvcGetEventVarNames(ioList []SrvcServiceIO, eventName string) []string {
	var result []string
	for _, io := range ioList {
		if io.EventName == eventName {
			result = append(result, io.VarName)
		}
	}
	return result
}

func isStringInArray(str string, arr []string) bool {
	for _, element := range arr {
		if element == str {
			return true
		}
	}
	return false
}
