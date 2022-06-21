package tinymap

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestTinyMap(t *testing.T) {
	var u TinyMap

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		u.SetBytes(key, i+5)
		testTinyMapGet(t, &u, key, i+5)
		u.SetBytes(key, i)
		testTinyMapGet(t, &u, key, i)
	}

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		testTinyMapGet(t, &u, key, i)
	}

	testTinyMapVisitValues(t, &u)

	u.Reset()

	testTinyMapVisitValues(t, &u)

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		testTinyMapGet(t, &u, key, nil)
	}

}

func testTinyMapGet(t *testing.T, u *TinyMap, key []byte, value interface{}) {
	v := u.GetBytes(key)
	if v == nil && value != nil {
		t.Fatalf("cannot obtain value for key=%q", key)
	}
	if !reflect.DeepEqual(v, value) {
		t.Fatalf("unexpected value for key=%q: %d. Expecting %d", key, v, value)
	}
}

func testTinyMapVisitValues(t *testing.T, u *TinyMap) {
	i := 0
	u.VisitValues(func(key []byte, val interface{}) {
		arr := *u
		if !bytes.Equal(key, arr[i].key) {
			t.Fatalf("unexpected key for item[%d]. Expecting %q, got %q", i, arr[i].key, key)
		}

		if !reflect.DeepEqual(arr[i].value, val) {
			t.Fatalf("unexpected value for key=%q: %v. Expecting %v", key, val, arr[i].value)
		}
		i++
	})
}

func TestTinyMapValueClose(t *testing.T) {
	var u TinyMap

	closeCalls := 0

	// store values implementing io.Closer
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		u.Set(key, &closerValue{&closeCalls})
	}

	// store values without io.Closer
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_noclose_%d", i)
		u.Set(key, i)
	}

	u.Reset()

	if closeCalls != 5 {
		t.Fatalf("unexpected number of Close calls: %d. Expecting 10", closeCalls)
	}
}

type closerValue struct {
	closeCalls *int
}

func (cv *closerValue) Close() error {
	(*cv.closeCalls)++
	return nil
}
