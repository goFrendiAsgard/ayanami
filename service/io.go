package service

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
