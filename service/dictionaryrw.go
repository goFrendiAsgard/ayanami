package service

import (
	"sync"
)

// DictionaryRW concurrent-safe setter/getter for Dictionary
type DictionaryRW struct {
	dictionary Dictionary
	lock       sync.RWMutex
}

// GetDictionary get DictionaryRW's dictionary
func (rw *DictionaryRW) GetDictionary() Dictionary {
	rw.lock.RLock()
	defer rw.lock.RUnlock()
	dictionary := rw.dictionary
	return dictionary
}

// Delete delete DictionaryRW's dictionary's key
func (rw *DictionaryRW) Delete(key string) {
	rw.lock.RLock()
	defer rw.lock.RUnlock()
	delete(rw.dictionary, key)
}

// Get get from dictionary
func (rw *DictionaryRW) Get(dottedKeys string) interface{} {
	rw.lock.RLock()
	defer rw.lock.RUnlock()
	return rw.dictionary.Get(dottedKeys)
}

// Has check whether key exists in dictionary
func (rw *DictionaryRW) Has(dottedKeys string) bool {
	rw.lock.RLock()
	defer rw.lock.RUnlock()
	return rw.dictionary.Has(dottedKeys)
}

// HasAll check whether all key in keyNames are available in dictionary
func (rw *DictionaryRW) HasAll(keyNames []string) bool {
	rw.lock.RLock()
	defer rw.lock.RUnlock()
	return rw.dictionary.HasAll(keyNames)
}

// Set set dictionary
func (rw *DictionaryRW) Set(dottedKeys string, newValue interface{}) error {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	return rw.dictionary.Set(dottedKeys, newValue)
}

// NewDictionaryRW create new DictionaryRW
func NewDictionaryRW() *DictionaryRW {
	return NewPresetDictionaryRW(make(Dictionary))
}

// NewPresetDictionaryRW create new DictionaryRW
func NewPresetDictionaryRW(dictionary Dictionary) *DictionaryRW {
	rw := DictionaryRW{dictionary: dictionary}
	return &rw
}
