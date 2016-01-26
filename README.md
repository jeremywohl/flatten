flatten
=======

[![GoDoc](https://godoc.org/github.com/jeremywohl/flatten?status.png)](https://godoc.org/github.com/jeremywohl/flatten)

Flatten makes flat, one-dimensional maps from arbitrarily nested ones, from JSON strings or Go native structures.  It can handles interior maps, slices and scalars.

Intended for JSON APIs, flatten operates on interior maps, slices and scalars.  Flat, compound
keys are generated in either dotted style (e.g. a.b.1.c) or Rails-like (e.g. a[b][1][c]).

You can flatten JSON strings.

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

Or Go maps directly.

    t := map[string]interface{}{
       "a": "b",
       "c": map[string]interface{}{
           "d": "e",
           "f": "g",
       },
       "z": 1.4567
    }
    
    flat, err := Flatten(nested, "", RAILS_STYLE)
    
    // output:
    // map[string]interface{}{
    //  "a":    "b",
    //  "c[d]": "e",
    //  "c[f]": "g",
    //  "z":    1.4567,
    // }

See [godoc](https://godoc.org/github.com/jeremywohl/flatten) for API.
