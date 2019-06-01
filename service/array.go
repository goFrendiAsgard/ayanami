package service

// IsStringInArray determine wheter a string is array's element or not
func IsStringInArray(str string, arr []string) bool {
	for _, element := range arr {
		if element == str {
			return true
		}
	}
	return false
}

// AppendUniqueString append string to array if string is not already in array
func AppendUniqueString(str string, arr []string) []string {
	if IsStringInArray(str, arr) {
		return arr
	}
	return append(arr, str)
}
