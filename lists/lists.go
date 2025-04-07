/*
Ordered map in simple key-value pairs.

NOTE: The code about [OrderedMap] is created before Go 1.18 that introduces generics.
I am going to keep the code compatible to old projects, but will still update if necessary.
If you find any problem or want to improve the code, please rise an issue or PR.
*/
package lists

import "errors"

// Predefined errors for the lists package.
var (
	ErrKeyAlreadyExists = errors.New("key already exists")
	ErrKeyNotFound      = errors.New("key not found")
	ErrPosOutOfRange    = errors.New("position out of range")
)

// Create a new, empty [OrderedMap].
func NewOrderedMap() *OrderedMap {
	omap := &OrderedMap{
		keys: make([]string, 0),
		data: make(map[string]interface{}),
	}
	return omap
}

// Ordered map that stores data in form of map[key]record in order of []key.
// It works like map[string]interface{} in mind of:
//   - A "key", the string in the map.
//   - A "record", the interface{} in the map.
//
// You may want to start with [NewOrderedMap] function.
type OrderedMap struct {
	keys []string               // to keep track of order of the data.
	data map[string]interface{} // for key-record paired data.
}

// Get all keys in order.
func (omap *OrderedMap) Keys() []string {
	return omap.keys
}

// Get all records in order of the keys.
func (omap *OrderedMap) Records() []interface{} {
	records := make([]interface{}, len(omap.keys))
	for i, key := range omap.keys {
		records[i] = omap.data[key]
	}
	return records
}

// Find a single record using a key.
func (omap *OrderedMap) Record(key string) (interface{}, error) {
	record, found := omap.data[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	return record, nil
}

// Insert a record at specified position.
// To add a record at the end without knowing length of keys, please use [OrderedMap.Append].
func (omap *OrderedMap) Insert(key string, record interface{}, position int) error {
	if position < 0 || position > len(omap.keys) {
		return ErrPosOutOfRange
	}
	if _, exists := omap.data[key]; exists {
		return ErrKeyAlreadyExists
	}
	// Process keys:
	newKeys := make([]string, len(omap.keys)+1)
	for i := range omap.keys {
		if i < position {
			newKeys[i] = omap.keys[i]
		} else {
			newKeys[i+1] = omap.keys[i]
		}
	}
	newKeys[position] = key
	omap.keys = newKeys
	// Process map data:
	omap.data[key] = record
	return nil
}

// Append to add a key-record pair at the end position.
func (omap *OrderedMap) Append(key string, record interface{}) error {
	if _, exists := omap.data[key]; exists {
		return ErrKeyAlreadyExists
	}
	omap.keys = append(omap.keys, key)
	omap.data[key] = record
	return nil
}

// Delete to remove a record using a key.
func (omap *OrderedMap) Delete(key string) error {
	// Delete requested record:
	if _, exists := omap.data[key]; !exists {
		return ErrKeyNotFound
	}
	delete(omap.data, key)
	// Delete the key from the slice:
	for i, k := range omap.keys {
		if k == key {
			omap.keys = append(omap.keys[:i], omap.keys[i+1:]...)
			break
		}
	}
	return nil
}

// Reset to clean existing order map.
func (omap *OrderedMap) Reset() {
	omap.keys = make([]string, 0)
	omap.data = make(map[string]interface{})
}

// Iterate over key-record pairs in order. This function returns a channel. And the channel yields structs containing "Key" and "Record" fields.
func (omap *OrderedMap) Iter() <-chan struct {
	Key    string
	Record interface{}
} {
	ch := make(chan struct {
		Key    string
		Record interface{}
	})
	go func() {
		defer close(ch)
		for _, key := range omap.keys {
			ch <- struct {
				Key    string
				Record interface{}
			}{key, omap.data[key]}
		}
	}()
	return ch
}
