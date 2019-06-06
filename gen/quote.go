package gen

import (
	"fmt"
	"strings"
)

// Quote add quote to a string
func Quote(str string) string {
	str = strings.ReplaceAll(str, "\"", "\\\"")
	return fmt.Sprintf("\"%s\"", str)
}

// QuoteArray quote elements of array and join it using delimiter
func QuoteArray(arr []string, delimiter string) string {
	for index, str := range arr {
		arr[index] = Quote(str)
	}
	return strings.Join(arr, delimiter)
}
