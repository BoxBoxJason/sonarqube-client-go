package generate

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/pkg/gojson"
)

// ConvertStringToStruct converts a JSON string to a Go struct definition.
func ConvertStringToStruct(json, name string) (string, error) {
	if json == "" {
		return "", errors.New("json string must not be empty")
	}

	reader := new(bytes.Buffer)
	reader.WriteString(json)

	byts, err := gojson.Generate(reader, gojson.ParseJSON, name, []string{"json"}, true, true)
	if err != nil {
		return "", fmt.Errorf("failed to generate struct: %w", err)
	}

	return string(byts), nil
}

// UnionJSONToStruct converts multiple JSON strings to a single Go struct definition.
func UnionJSONToStruct(jsons []string, name string) (string, error) {
	if len(jsons) == 0 {
		return "", errors.New("jsons string must not be zero")
	}

	reader := new(bytes.Buffer)
	reader.WriteString("[")
	reader.WriteString(strings.Join(jsons, ","))
	reader.WriteString("]")

	byts, err := gojson.Generate(reader, gojson.ParseJSON, name, []string{"json"}, false, true)
	if err != nil {
		return "", fmt.Errorf("failed to generate struct: %w", err)
	}

	return strings.Replace(string(byts), "[]", "", 1), nil
}
