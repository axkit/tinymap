// The package tinymap provides functionality of map as slice of key-value pair.
//
// The type TinyMap is copy of "userData" type from github.com/valyala/fasthttp
//
// All credits go to Aliaksandr Valialkin, VertaMedia, Kirill Danshin, Erik Dubbelboer, FastHTTP Authors
package tinymap

import (
	"io"
	"unsafe"
)

// item represents a key-value pair where key is a byte slice and value is an interface{}
type item struct {
	key   []byte
	value interface{}
}

// TinyMap is a slice of items, acting as a map.
type TinyMap []item

// Set adds or updates a key-value pair in the TinyMap.
// If the key exists, its value is updated. If not, the key-value pair is added.
func (tm *TinyMap) Set(key string, value interface{}) {
	args := *tm
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		// Check if the key already exists
		if string(kv.key) == key {
			kv.value = value
			return
		}
	}

	// Check if there is capacity to add a new item without reallocating
	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		// Copy the key and assign the value
		kv.key = append(kv.key[:0], key...)
		kv.value = value
		*tm = args
		return
	}

	// Add a new key-value pair when capacity needs to grow
	kv := item{}
	kv.key = append(kv.key[:0], key...)
	kv.value = value
	*tm = append(args, kv)
}

// SetBytes adds or updates a key-value pair using a byte slice as the key.
// It internally converts the byte slice to a string before setting the value.
func (tm *TinyMap) SetBytes(key []byte, value interface{}) {
	tm.Set(b2s(key), value)
}

// Get retrieves the value associated with a key.
// Returns nil if the key does not exist.
func (tm *TinyMap) Get(key string) interface{} {
	args := *tm
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		// Check if the key matches
		if string(kv.key) == key {
			return kv.value
		}
	}
	return nil
}

// GetBytes retrieves the value associated with a byte slice key.
// It internally converts the byte slice to a string before getting the value.
func (tm *TinyMap) GetBytes(key []byte) interface{} {
	return tm.Get(b2s(key))
}

// Reset clears the TinyMap by truncating it to zero length.
// If any value implements io.Closer, it will call the Close method before clearing.
func (tm *TinyMap) Reset() {
	args := *tm
	n := len(args)
	for i := 0; i < n; i++ {
		v := args[i].value
		// Check if value implements io.Closer and close it
		if vc, ok := v.(io.Closer); ok {
			vc.Close()
		}
	}
	*tm = (*tm)[:0]
}

// b2s converts a byte slice to a string without allocation.
// Unsafe operation; use with caution.
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// VisitValues applies a visitor function to all key-value pairs in the TinyMap.
// The visitor function receives the key as a byte slice and the value.
func (tm *TinyMap) VisitValues(visitor func([]byte, interface{})) {
	arr := *tm
	for i, n := 0, len(arr); i < n; i++ {
		kv := &arr[i]
		// Call the visitor function for each key-value pair
		visitor(kv.key, kv.value)
	}
}
