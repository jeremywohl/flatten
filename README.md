flatten
=======

[![GoDoc](https://godoc.org/github.com/jeremywohl/flatten?status.png)](https://godoc.org/github.com/jeremywohl/flatten)
[![Build Status](https://travis-ci.org/jeremywohl/flatten.svg?branch=master)](https://travis-ci.org/jeremywohl/flatten)

Flatten makes flat, one-dimensional maps from arbitrarily nested ones.

Map keys turn into compound
names, like `a.b.1.c` (dotted style), `a[b][1][c]` (Rails style) or `a/b/1/c` (path style).  It takes input as either JSON strings or
Go structures.  It knows how to traverse JSON types: maps, slices and scalars.

You can flatten JSON strings.

```go
nested := `{
  "one": {
    "two": [
      "2a",
      "2b"
    ]
  },
  "side": "value"
}`

flat, err := FlattenString(nested, "", DOT_STYLE)

// output: `{ "one.two.0": "2a", "one.two.1": "2b", "side": "value" }`
```

Or Go maps directly.

```go
t := map[string]interface{}{
   "a": "b",
   "c": map[string]interface{}{
       "d": "e",
       "f": "g",
   },
   "z": 1.4567,
}

flat, err := Flatten(nested, "", RAILS_STYLE)

// output:
// map[string]interface{}{
//  "a":    "b",
//  "c[d]": "e",
//  "c[f]": "g",
//  "z":    1.4567,
// }
```

See [godoc](https://godoc.org/github.com/jeremywohl/flatten) for API.
