package service

import (
	"reflect"
	"strconv"
	"strings"
)

// Dictionary map[string]interface{}
type Dictionary map[string]interface{}

// Get get from dictionary
func (dictionary Dictionary) Get(dottedKeys string) interface{} {
	keyParts := strings.Split(dottedKeys, ".")
	var data interface{}
	data = dictionary
	for _, key := range keyParts {
		reflectInfo := reflect.ValueOf(data)
		reflectKind := reflectInfo.Kind()
		if reflectKind == reflect.Map {
			var exists bool
			if data, exists = data.(Dictionary)[key]; !exists {
				return nil
			}
		} else if reflectKind == reflect.Slice {
			index, err := strconv.ParseInt(key, 10, 64)
			if err != nil {
				return nil
			}
			var arr []interface{}
			arr = data.([]interface{})
			if int64(len(arr)) <= index {
				return nil
			}
			data = arr[index]
		}
	}
	return data
}
