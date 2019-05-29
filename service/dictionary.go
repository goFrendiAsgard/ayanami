package service

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// aliases
type arrayOfInterface = []interface{}
type mapOfInterface = map[string]interface{}

// Dictionary map[string]interface{}
type Dictionary mapOfInterface

// Get get from dictionary
func (dictionary Dictionary) Get(dottedKeys string) interface{} {
	keyParts := strings.Split(dottedKeys, ".")
	var data interface{}
	data = dictionary
	for _, key := range keyParts {
		reflectKind := getReflectKind(data)
		if reflectKind == reflect.Map {
			var exists bool
			if data, exists = getDictionary(data)[key]; !exists {
				return nil
			}
		} else if reflectKind == reflect.Slice {
			index, err := getArrayIndex(key)
			if err != nil {
				return nil
			}
			var arr arrayOfInterface
			arr = data.(arrayOfInterface)
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
		reflectKind := getReflectKind(data)
		if reflectKind == reflect.Map {
			var exists bool
			if data, exists = getDictionary(data)[key]; !exists {
				return false
			}
		} else if reflectKind == reflect.Slice {
			index, err := getArrayIndex(key)
			if err != nil {
				return false
			}
			var arr arrayOfInterface
			arr = data.(arrayOfInterface)
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
	var pointer interface{} = *dictionary
	keyParts := strings.Split(dottedKeys, ".")
	if len(keyParts) == 1 { // dottedKeys is not nested, set the value and return
		pointer.(Dictionary)[dottedKeys] = newValue
		return nil
	}
	// dottedKeys is nested
	return setDictionaryChildren(pointer, keyParts, newValue)
}

func setDictionaryChildren(pointer interface{}, keyParts []string, newValue interface{}) error {
	// dottedKeys is nested
	lastKeyIndex := len(keyParts) - 1
	for keyIndex, key := range keyParts {
		subKeyIndex := keyIndex + 1
		isLastKey := keyIndex == lastKeyIndex
		reflectKind := getReflectKind(pointer)
		if reflectKind == reflect.Map { // pointer refer to map
			if isLastKey { // this is the last key, set value and exit
				pointer.(Dictionary)[key] = newValue
				return nil
			}
			if pointer.(Dictionary).Has(key) {
				subKey := keyParts[subKeyIndex]
				subVal := pointer.(Dictionary)[key]
				reflectSubValKind := getReflectKind(subVal)
				if reflectSubValKind == reflect.Slice { // child is array
					arrayIndex, err := getArrayIndex(subKey)
					if err != nil {
						return err
					}
					// fill up child with empty dictionaries
					pointer.(Dictionary)[key] = fillArray(subVal.(arrayOfInterface), arrayIndex)
				} else if reflectSubValKind == reflect.Map {
					// make sure every map is converted into Dictionary
					pointer.(Dictionary)[key] = getDictionary(pointer.(Dictionary)[key])
				}
			} else {
				pointer.(Dictionary)[key] = make(Dictionary)
			}
			pointer = pointer.(Dictionary)[key]
		} else if reflectKind == reflect.Slice { // pointer refer to array
			arrayIndex, _ := getArrayIndex(key)
			if isLastKey { // this is the last key, set value and exit
				pointer.(arrayOfInterface)[arrayIndex] = newValue
				return nil
			}
			subVal := pointer.(arrayOfInterface)[arrayIndex]
			reflectSubValKind := getReflectKind(subVal)
			if reflectSubValKind == reflect.Map {
				pointer.(arrayOfInterface)[arrayIndex] = getDictionary(pointer.(arrayOfInterface)[arrayIndex])
			}
			pointer = pointer.(arrayOfInterface)[arrayIndex]
		} else {
			return errors.New("cannot override non-dictionary and non-list")
		}
	}
	return nil
}

func getArrayIndex(key string) (int64, error) {
	index, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, err
	}
	if index < 0 {
		return 0, errors.New("negative array index is not allowed")
	}
	return index, nil
}

func getDictionary(data interface{}) Dictionary {
	mapVal, isMapVal := data.(mapOfInterface)
	if isMapVal {
		return Dictionary(mapVal)
	}
	return data.(Dictionary)
}

func getReflectKind(data interface{}) reflect.Kind {
	reflectInfo := reflect.ValueOf(data)
	return reflectInfo.Kind()
}

func fillArray(array arrayOfInterface, newMaxIndex int64) arrayOfInterface {
	arrayLength := int64(len(array))
	for index := arrayLength; index <= newMaxIndex; index++ {
		array = append(array, make(Dictionary))
	}
	return array
}
