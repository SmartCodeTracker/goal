package run

import (
	"fmt"

	"github.com/kylelemons/go-gypsy/yaml"
)

// parseMap gets a yaml.Map and a key which's value must be transformed
// into map[string][]string.
// It returns an error if the transformation is not possible.
//
// It is used for parsing the following kind of structure:
//	watch:
//		"./controllers/*":
//			- go build ./
//			- go test ./
func parseMap(m yaml.Map, key string) (map[string][]string, error) {
	// Make sure the requested key is presented.
	t, ok := m[key]
	if !ok {
		return nil, fmt.Errorf(`no "%s" key found in YAML configuration file`, key)
	}

	// Make sure the value is a map.
	val, ok := t.(yaml.Map)
	if !ok {
		return nil, fmt.Errorf(`"%s" must be a map`, key)
	}

	// Make sure every element of the map is of the expected []string type.
	res := make(map[string][]string, len(val))
	for k := range val {
		lst, err := parseSlice(val, k)
		if err != nil {
			return nil, err
		}
		res[k] = lst
	}
	return res, nil
}

// parseSlice gets a yaml.Map and a key which's value should be transformed into
// a slice of strings.
// It returns an error if the transformation is not possible.
//
// It is used for parsing the following kind of structure:
//	init:
//		- go build ./
//		- go test ./
func parseSlice(m yaml.Map, key string) ([]string, error) {
	// Make sure the requested key is presented.
	ts, ok := m[key]
	if !ok {
		return nil, fmt.Errorf(`no "%s" key found in YAML configuration file`, key)
	}

	// Check whether values are of List type.
	lst, ok := ts.(yaml.List)
	if !ok {
		return nil, fmt.Errorf(`"%s" must be a list`, key)
	}

	// Make sure every element of the list is a string.
	res := make([]string, len(lst))
	for i := 0; i < len(lst); i++ {
		// Make sure the element is a string.
		s, ok := lst[i].(yaml.Scalar)
		if !ok {
			return nil, fmt.Errorf(`%d-th element of the "%s" list is not a string`, i, key)
		}

		// If so, add the element to the list.
		res = append(res, s.String())
	}
	return res, nil
}
