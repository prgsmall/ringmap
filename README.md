# 🔃 github.com/prgsmall/ringmap [![GoDoc](https://godoc.org/github.com/prgsmall/ringmap?status.svg)](https://godoc.org/github.com/prgsmall/ringmap) [![Build Status](https://travis-ci.org/elliotchance/ringmap.svg?branch=master)](https://travis-ci.org/elliotchance/ringmap)

## Installation

```bash
go get -u github.com/prgsmall/ringmap
```

This wraps the orderedmap data structure available here:  https://github.com/elliotchance/orderedmap

## Basic Usage

A `*RingMap` is a high performance ordered map that maintains amortized O(1)
for `Set`, `Get`, `Delete` and `Len`:

```go
m := ringmap.NewRingMap()

m.Set("foo", "bar")
m.Set("qux", 1.23)
m.Set(123, true)

m.Delete("qux")

m.Put("zzz", "yyy") // Deletes if the key exists, then calls Set
```

Internally an `*RingMap` uses a combination of a map and linked list.

## Iterating

Be careful using `Keys()` as it will create a copy of all of the keys so it's
only suitable for a small number of items:

```go
for _, key := range m.Keys() {
	value, _:= m.Get(key)
	fmt.Println(key, value)
}
```

For larger maps you should use `Front()` or `Back()` to iterate per element:

```go
// Iterate through all elements from oldest to newest:
for el := m.Front(); el != nil; el = el.Next() {
    fmt.Println(el.Key, el.Value)
}

// You can also use Back and Prev to iterate in reverse:
for el := m.Back(); el != nil; el = el.Prev() {
    fmt.Println(el.Key, el.Value)
}
```

The iterator is safe to use bidirectionally, and will return `nil` once it goes
beyond the first or last item.

If the map is changing while the iteration is in-flight it may produce
unexpected behavior.
