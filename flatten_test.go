package flatten

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		test   string
		want   map[string]interface{}
		prefix string
		style  SeparatorStyle
	}{
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			map[string]interface{}{
				"foo.jim":      "bean",
				"fee":          "bar",
				"n1.alist.0":   "a",
				"n1.alist.1":   "b",
				"n1.alist.2":   "c",
				"n1.alist.3.d": "other",
				"n1.alist.3.e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			"",
			DotStyle,
		},
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
					"a",
					"b",
					"c",
					{
						"d": "other",
						"e": "another"
					}
					]
				}
			}`,
			map[string]interface{}{
				"foo[jim]":        "bean",
				"fee":             "bar",
				"n1[alist][0]":    "a",
				"n1[alist][1]":    "b",
				"n1[alist][2]":    "c",
				"n1[alist][3][d]": "other",
				"n1[alist][3][e]": "another",
			},
			"",
			RailsStyle,
		},
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			map[string]interface{}{
				"foo/jim":      "bean",
				"fee":          "bar",
				"n1/alist/0":   "a",
				"n1/alist/1":   "b",
				"n1/alist/2":   "c",
				"n1/alist/3/d": "other",
				"n1/alist/3/e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			"",
			PathStyle,
		},
		{
			`{ "a": { "b": "c" }, "e": "f" }`,
			map[string]interface{}{
				"p:a.b": "c",
				"p:e":   "f",
			},
			"p:",
			DotStyle,
		},
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			map[string]interface{}{
				"foo_jim":      "bean",
				"fee":          "bar",
				"n1_alist_0":   "a",
				"n1_alist_1":   "b",
				"n1_alist_2":   "c",
				"n1_alist_3_d": "other",
				"n1_alist_3_e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			"",
			UnderscoreStyle,
		},
	}

	for i, test := range cases {
		var m interface{}
		err := json.Unmarshal([]byte(test.test), &m)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
			continue
		}
		got, err := Flatten(m.(map[string]interface{}), test.prefix, test.style)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v wanted: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenString(t *testing.T) {
	cases := []struct {
		test   string
		want   string
		prefix string
		style  SeparatorStyle
		err    error
	}{
		// 1
		{
			`{ "a": "b" }`,
			`{ "a": "b" }`,
			"",
			DotStyle,
			nil,
		},
		// 2
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			`{ "a.b.c.d": "e", "bool": true, "number": 1.4567 }`,
			"",
			DotStyle,
			nil,
		},
		// 3
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			`{ "a/b/c/d": "e", "bool": true, "number": 1.4567 }`,
			"",
			PathStyle,
			nil,
		},
		// 4
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			`{ "a--b--c--d": "e" }`,
			"",
			SeparatorStyle{Middle: "--"}, // emdash
			nil,
		},
		// 5
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			`{ "a(b)(c)(d)": "e" }`,
			"",
			SeparatorStyle{Before: "(", After: ")"}, // paren groupings
			nil,
		},
		// 6 -- with leading whitespace
		{
			`
			  	{ "a": { "b" : { "c" : { "d" : "e" } } } }`,
			`{ "a(b)(c)(d)": "e" }`,
			"",
			SeparatorStyle{Before: "(", After: ")"}, // paren groupings
			nil,
		},

		//
		// Valid JSON text, but invalid for FlattenString
		//

		// 7
		{
			`[ "a": { "b": "c" }, "d" ]`,
			`bogus`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
		// 8
		{
			``,
			`bogus`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
		// 9
		{
			`astring`,
			`bogus`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
		// 10
		{
			`false`,
			`bogus`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
		// 11
		{
			`42`,
			`bogus`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
		// 12 -- prior to version 1.0.1, this was accepted & unmarshalled as an empty map, finally returning `{}`.
		{
			`null`,
			`{}`,
			"",
			PathStyle,
			NotValidJsonInputError,
		},
	}

	for i, test := range cases {
		got, err := FlattenString(test.test, test.prefix, test.style)
		if err != test.err {
			t.Errorf("%d: error mismatch, got: [%v], wanted: [%v]", i+1, err, test.err)
			continue
		}
		if err != nil {
			continue
		}

		nixws := func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}

		if got != strings.Map(nixws, test.want) {
			t.Errorf("%d: mismatch, got: %v wanted: %v", i+1, got, test.want)
		}
	}
}
