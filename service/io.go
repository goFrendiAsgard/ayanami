package service

// IO single IO configuration
type IO struct {
	VarName   string
	EventName string
}

// GetUniqueVarNames get unique varNames from IO list
func GetUniqueVarNames(ioList []IO) []string {
	var result []string
	for _, io := range ioList {
		if !isStringInArray(io.VarName, result) {
			result = append(result, io.VarName)
		}
	}
	return result
}

// GetVarEventNames get eventNames from IO list with specified varName
func GetVarEventNames(ioList []IO, varName string) []string {
	var result []string
	for _, io := range ioList {
		if io.VarName == varName {
			result = append(result, io.EventName)
		}
	}
	return result
}

// GetUniqueEventNames get unique eventNames from IO list
func GetUniqueEventNames(ioList []IO) []string {
	var result []string
	for _, io := range ioList {
		if !isStringInArray(io.EventName, result) {
			result = append(result, io.EventName)
		}
	}
	return result
}

// GetEventVarNames get varNames from IO list with specified eventName
func GetEventVarNames(ioList []IO, eventName string) []string {
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
