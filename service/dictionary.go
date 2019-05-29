package service

import (
	"errors"
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
			index, err := keyToIndex(key)
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

// Has is key exists
func (dictionary Dictionary) Has(dottedKeys string) bool {
	keyParts := strings.Split(dottedKeys, ".")
	var data interface{}
	data = dictionary
	for _, key := range keyParts {
		reflectInfo := reflect.ValueOf(data)
		reflectKind := reflectInfo.Kind()
		if reflectKind == reflect.Map {
			var exists bool
			if data, exists = data.(Dictionary)[key]; !exists {
				return false
			}
		} else if reflectKind == reflect.Slice {
			index, err := keyToIndex(key)
			if err != nil {
				return false
			}
			var arr []interface{}
			arr = data.([]interface{})
			if int64(len(arr)) <= index {
				return false
			}
			data = arr[index]
		}
	}
	return true
}

// HasAll check whether all key in keyNames are available in dictionary
func (dictionary *Dictionary) HasAll(keyNames []string) bool {
	for _, keyName := range keyNames {
		if !dictionary.Has(keyName) {
			return false
		}
	}
	return true
}

// Set set dictionary
func (dictionary *Dictionary) Set(dottedKeys string, newValue interface{}) error {
	var pointer interface{}
	pointer = *dictionary
	keyParts := strings.Split(dottedKeys, ".")
	if len(keyParts) == 1 {
		pointer.(Dictionary)[dottedKeys] = newValue
		return nil
	}
	lastKeyIndex := len(keyParts) - 1
	for keyIndex, key := range keyParts {
		reflectInfo := reflect.ValueOf(pointer)
		reflectKind := reflectInfo.Kind()
		if reflectKind == reflect.Map {
			_, exists := pointer.(Dictionary)[key]
			if exists && keyIndex < lastKeyIndex {
				subKey := keyParts[keyIndex+1]
				subVal := pointer.(Dictionary)[key]
				reflectSubValInfo := reflect.ValueOf(subVal)
				reflectSubValKind := reflectSubValInfo.Kind()
				if reflectSubValKind == reflect.Slice {
					index, err := keyToIndex(subKey)
					if err != nil {
						return err
					}
					arrayLength := int64(len(pointer.(Dictionary)[key].([]interface{})))
					// if array length is less than assigned value, populate the array until it has sensible length
					if index >= arrayLength {
						for i := arrayLength; i <= index; i++ {
							pointer.(Dictionary)[key] = append(pointer.(Dictionary)[key].([]interface{}), "")
						}
					}
					// assign array
					if keyIndex+1 == lastKeyIndex {
						pointer.(Dictionary)[key].([]interface{})[index] = newValue
					} else {
						pointer.(Dictionary)[key].([]interface{})[index] = make(Dictionary)
					}
				}
			} else if !exists && keyIndex < lastKeyIndex {
				pointer.(Dictionary)[key] = make(Dictionary)
			} else if keyIndex == lastKeyIndex {
				pointer.(Dictionary)[key] = newValue
			}
			pointer = pointer.(Dictionary)[key]
		} else if reflectKind == reflect.Slice {
			index, _ := keyToIndex(key)
			pointer = pointer.([]interface{})[index]
		} else {
			return errors.New("cannot override non-dictionary and non-list")
		}
	}
	return nil
}

func keyToIndex(key string) (int64, error) {
	index, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, err
	}
	if index < 0 {
		return 0, errors.New("negative array index is not allowed")
	}
	return index, nil
}
