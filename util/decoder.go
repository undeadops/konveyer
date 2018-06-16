package util

import (
	"bytes"
	"unicode"

	"github.com/ghodss/yaml"
)

// ToJSON converts a single YAML document into a JSON document
// or returns an error. If the document appears to be JSON the
// YAML decoding path is not used.
func ToJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml.YAMLToJSON(data)
}

var jsonPrefix = []byte("{")

func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, jsonPrefix)
}

func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}
