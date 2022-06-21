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

type item struct {
	key   []byte
	value interface{}
}

type TinyMap []item

func (d *TinyMap) Set(key string, value interface{}) {
	args := *d
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if string(kv.key) == key {
			kv.value = value
			return
		}
	}

	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		kv.key = append(kv.key[:0], key...)
		kv.value = value
		*d = args
		return
	}

	kv := item{}
	kv.key = append(kv.key[:0], key...)
	kv.value = value
	*d = append(args, kv)
}

func (d *TinyMap) SetBytes(key []byte, value interface{}) {
	d.Set(b2s(key), value)
}

func (d *TinyMap) Get(key string) interface{} {
	args := *d
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if string(kv.key) == key {
			return kv.value
		}
	}
	return nil
}

func (d *TinyMap) GetBytes(key []byte) interface{} {
	return d.Get(b2s(key))
}

func (d *TinyMap) Reset() {
	args := *d
	n := len(args)
	for i := 0; i < n; i++ {
		v := args[i].value
		if vc, ok := v.(io.Closer); ok {
			vc.Close()
		}
	}
	*d = (*d)[:0]
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (d *TinyMap) VisitValues(visitor func([]byte, interface{})) {
	arr := *d
	for i, n := 0, len(arr); i < n; i++ {
		kv := &arr[i]
		visitor(kv.key, kv.value)
	}
}
