package flatten

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

type KeyMerger interface {
	MergeKeys(top bool, key, subkey string) string
}

type FuncMerger func(top bool, key, subkey string) string

func (f FuncMerger) MergeKeys(top bool, key, subkey string) string {
	return f(top, key, subkey)
}

// The style of keys.  If there is an input with two
// nested keys "f" and "g", with "f" at the root,
//    { "f": { "g": ... } }
// the output will be the concatenation
//    f{Middle}{Before}g{After}...
// Any struct element may be blank.
// If you use Middle, you will probably leave Before & After blank, and vice-versa.
// See examples in flatten_test.go and the "Default styles" here.
type SeparatorStyle struct {
	Before string // Prepend to key
	Middle string // Add between keys
	After  string // Append to key
}

func (style SeparatorStyle) MergeKeys(top bool, key, subkey string) string {
	if top {
		return key + subkey
	}

	return key + style.Before + style.Middle + subkey + style.After
}

// Default styles
var (
	// Separate nested key components with dots, e.g. "a.b.1.c.d"
	DotStyle = SeparatorStyle{Middle: "."}

	// Separate with path-like slashes, e.g. a/b/1/c/d
	PathStyle = SeparatorStyle{Middle: "/"}

	// Separate ala Rails, e.g. "a[b][c][1][d]"
	RailsStyle = SeparatorStyle{Before: "[", After: "]"}

	// Separate with underscores, e.g. "a_b_1_c_d"
	UnderscoreStyle = SeparatorStyle{Middle: "_"}
)

// Nested input must be a map or slice
var NotValidInputError = errors.New("Not a valid input: map or slice")

// Flatten generates a flat map from a nested one.  The original may include values of type map, slice and scalar,
// but not struct.  Keys in the flat map will be a compound of descending map keys and slice iterations.
// The presentation of keys is set by style.  A prefix is joined to each key.
func Flatten(nested map[string]interface{}, prefix string, keyMerger KeyMerger) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, keyMerger)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

// JSON nested input must be a map
var NotValidJsonInputError = errors.New("Not a valid input, must be a map")

var isJsonMap = regexp.MustCompile(`^\s*\{`)

// FlattenString generates a flat JSON map from a nested one.  Keys in the flat map will be a compound of
// descending map keys and slice iterations.  The presentation of keys is set by style.  A prefix is joined
// to each key.
func FlattenString(nestedstr, prefix string, keyMerger KeyMerger) (string, error) {
	if !isJsonMap.MatchString(nestedstr) {
		return "", NotValidJsonInputError
	}

	var nested map[string]interface{}
	err := json.Unmarshal([]byte(nestedstr), &nested)
	if err != nil {
		return "", err
	}

	flatmap, err := Flatten(nested, prefix, keyMerger)
	if err != nil {
		return "", err
	}

	flatb, err := json.Marshal(&flatmap)
	if err != nil {
		return "", err
	}

	return string(flatb), nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, keyMerger KeyMerger) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey, keyMerger); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := keyMerger.MergeKeys(top, prefix, k)
			assign(newKey, v)
		}
	case []interface{}:
		for i, v := range nested.([]interface{}) {
			newKey := keyMerger.MergeKeys(top, prefix, strconv.Itoa(i))
			assign(newKey, v)
		}
	default:
		return NotValidInputError
	}

	return nil
}
