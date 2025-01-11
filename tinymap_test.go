package tinymap

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestTinyMap(t *testing.T) {
	var tm TinyMap

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		tm.SetBytes(key, i+5)
		testTinyMapGet(t, &tm, key, i+5)
		tm.SetBytes(key, i)
		testTinyMapGet(t, &tm, key, i)
	}

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		testTinyMapGet(t, &tm, key, i)
	}

	testTinyMapVisitValues(t, &tm)

	tm.Reset()

	testTinyMapVisitValues(t, &tm)

	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		testTinyMapGet(t, &tm, key, nil)
	}

}

func testTinyMapGet(t *testing.T, tm *TinyMap, key []byte, value interface{}) {
	v := tm.GetBytes(key)
	if v == nil && value != nil {
		t.Fatalf("cannot obtain value for key=%q", key)
	}
	if !reflect.DeepEqual(v, value) {
		t.Fatalf("unexpected value for key=%q: %d. Expecting %d", key, v, value)
	}
}

func testTinyMapVisitValues(t *testing.T, tm *TinyMap) {
	i := 0
	tm.VisitValues(func(key []byte, val interface{}) {
		arr := *tm
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
	var tm TinyMap

	closeCalls := 0

	// store values implementing io.Closer
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key_%d", i)
		tm.Set(key, &closerValue{&closeCalls})
	}

	// store values without io.Closer
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key_noclose_%d", i)
		tm.Set(key, i)
	}

	tm.Reset()

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
