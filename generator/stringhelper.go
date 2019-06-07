package generator

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// StringHelper string helper
type StringHelper struct{}

// Quote add quote to a string
func (s StringHelper) Quote(str string) string {
	str = strings.ReplaceAll(str, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", str)
}

// QuoteMap quote array's elements
func (s StringHelper) QuoteMap(data map[string]string) map[string]string {
	for index, str := range data {
		data[index] = s.Quote(str)
	}
	return data
}

// QuoteArray quote array's elements
func (s StringHelper) QuoteArray(arr []string) []string {
	for index, str := range arr {
		arr[index] = s.Quote(str)
	}
	return arr
}

// QuoteArrayAndJoin quote elements of array and join it using delimiter
func (s StringHelper) QuoteArrayAndJoin(arr []string, delimiter string) string {
	arr = s.QuoteArray(arr)
	return strings.Join(arr, delimiter)
}

// IsMatch check whether is string is match with regex pattern or not
func (s StringHelper) IsMatch(str, regexPattern string) bool {
	reg, err := regexp.Compile(regexPattern)
	if err != nil {
		log.Printf("[ERROR] Invalid regex expression `%s`", regexPattern)
		return false
	}
	return reg.Match([]byte(str))
}

// IsAlpha check whether string contains letter only or not
func (s StringHelper) IsAlpha(str string) bool {
	return s.IsMatch(str, `^[a-zA-Z]+$`)
}

// IsNumeric check whether string contains numbers only or not
func (s StringHelper) IsNumeric(str string) bool {
	return s.IsMatch(str, `^[0-9]+$`)
}

// IsAlphaNumeric check whether string contains alphanumeric or not
func (s StringHelper) IsAlphaNumeric(str string) bool {
	return s.IsMatch(str, `^[a-zA-Z0-9]+$`)
}
