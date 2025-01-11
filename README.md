# tinymap

[![Build Status](https://github.com/axkit/bitset/actions/workflows/go.yml/badge.svg)](https://github.com/axkit/tinymap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axkit/tinymap)](https://goreportcard.com/report/github.com/axkit/tinymap)
[![GoDoc](https://pkg.go.dev/badge/github.com/axkit/tinymap)](https://pkg.go.dev/github.com/axkit/tinymap)
[![Coverage Status](https://coveralls.io/repos/github/axkit/tinymap/badge.svg?branch=main)](https://coveralls.io/github/axkit/tinymap?branch=main)


`tinymap` is a Go package providing a slice-based map implementation that stores key-value pairs as a slice of structs. It is inspired by and based on the `userData` type from [github.com/valyala/fasthttp](https://github.com/valyala/fasthttp).

This package is particularly useful in scenarios where you need a lightweight and straightforward map-like structure with a predictable memory layout.

## Features

- **Slice-based Map**: Keys and values are stored as slices of key-value pairs.
- **Custom Key Types**: Supports string and `[]byte` keys.
- **Memory-Efficient Reset**: Resets the map while ensuring proper cleanup of resources implementing the `io.Closer` interface.

## Installation

Install the package using `go get`:

```sh
$ go get github.com/axkit/tinymap
```

## Usage

Here is how you can use the `tinymap` package:

### Basic Operations

```go
package main

import (
    "fmt"
    "github.com/axkit/tinymap"
)

func main() {
    var tm tinymap.TinyMap

    // Set key-value pairs
    tm.Set("foo", 42)
    tm.Set("bar", "hello")

    // Retrieve values by key
    fmt.Println(tm.Get("foo"))  // Output: 42
    fmt.Println(tm.Get("bar"))  // Output: hello

    // Use byte slice keys
    tm.SetBytes([]byte("baz"), 3.14)
    fmt.Println(tm.GetBytes([]byte("baz")))  // Output: 3.14

    // Visit all key-value pairs
    tm.VisitValues(func(key []byte, value interface{}) {
        fmt.Printf("Key: %s, Value: %v\n", key, value)
    })

    // Reset the map
    tm.Reset()
    fmt.Println(tm.Get("foo")) // Output: <nil>
}
```

### Benchmarking

You can compare the performance of `tinymap` against standard Go maps by running the provided benchmarks:

```sh
$ go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/axkit/tinymap
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
BenchmarkTinyMapCustom-8        53070952                20.79 ns/op            0 B/op          0 allocs/op
BenchmarkTinyMapStdMap-8        30541206                37.11 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/axkit/tinymap        2.306s
```


## API Reference

### `TinyMap`

#### Methods

- **`Set(key string, value interface{})`**: Adds or updates a key-value pair.
- **`SetBytes(key []byte, value interface{})`**: Adds or updates a key-value pair using a byte slice as the key.
- **`Get(key string) interface{}`**: Retrieves the value for a given key. Returns `nil` if the key does not exist.
- **`GetBytes(key []byte) interface{}`**: Retrieves the value for a byte slice key.
- **`Reset()`**: Clears all key-value pairs, calling `Close` on any values that implement the `io.Closer` interface.
- **`VisitValues(visitor func([]byte, interface{}))`**: Iterates through all key-value pairs, calling the visitor function with each pair.

## Testing

Run the tests to ensure the package works as expected:

```sh
$ go test ./...
```

## Examples

### Advanced Example

```go
package main

import (
    "fmt"
    "io"
    "github.com/axkit/tinymap"
)

type resource struct {
    name string
}

func (r *resource) Close() error {
    fmt.Printf("Closing resource: %s\n", r.name)
    return nil
}

func main() {
    var tm tinymap.TinyMap

    // Add resources that implement io.Closer
    tm.Set("res1", &resource{name: "Resource1"})
    tm.Set("res2", &resource{name: "Resource2"})

    // Reset and observe resource cleanup
    tm.Reset() // Output: Closing resource: Resource1
               //         Closing resource: Resource2
}
```

## License

This package is open-sourced under the MIT license. See the LICENSE file for details.
