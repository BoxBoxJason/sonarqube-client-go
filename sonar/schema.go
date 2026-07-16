package sonar

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"strings"
)

// SchemaMismatch describes a JSON object field observed in a decoded API
// response that has no corresponding field on the destination Go struct. It
// signals that a modeled response type has drifted from what the SonarQube
// API actually returns: either a field is missing from the struct, or the
// struct models it under the wrong name.
type SchemaMismatch struct {
	// Endpoint identifies the request the mismatch was observed on, formatted
	// as "METHOD scheme://host/path".
	Endpoint string
	// GoType is the name of the struct type being validated.
	GoType string
	// Path is the dotted/bracketed path to the unmapped field, e.g.
	// "components[].measures[].bestValue".
	Path string
}

// String formats the mismatch as a single human-readable line.
func (m SchemaMismatch) String() string {
	return fmt.Sprintf("%s: field %q has no match in %s", m.Endpoint, m.Path, m.GoType)
}

// SchemaObserver is invoked by Client.Do, when configured via
// WithSchemaObserver, after every successfully decoded API response. It
// receives any SchemaMismatch found between the raw response body and the
// destination struct for that call. mismatches is empty when the response
// matched its destination type exactly.
type SchemaObserver func(endpoint string, mismatches []SchemaMismatch)

// CheckSchema compares a raw JSON response body against the exported,
// JSON-tagged fields declared on the type of dest and returns one
// SchemaMismatch per JSON object key that has no corresponding struct field.
// Object keys are only checked where the JSON value lines up against a Go
// struct; JSON matched against a Go map is left unchecked since maps
// intentionally accept arbitrary keys.
//
// CheckSchema returns an error only if data is not valid JSON. A nil dest, or
// a dest whose underlying type is never a struct anywhere in its shape,
// simply yields no mismatches.
func CheckSchema(endpoint string, data []byte, dest any) ([]SchemaMismatch, error) {
	if dest == nil {
		return nil, nil
	}

	var raw any

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body for schema check: %w", err)
	}

	destType := indirectType(reflect.TypeOf(dest))
	if destType == nil {
		return nil, nil
	}

	var mismatches []SchemaMismatch

	walkSchema(raw, destType, "", endpoint, &mismatches)

	return mismatches, nil
}

// indirectType dereferences pointer types until it reaches the underlying type.
func indirectType(fieldType reflect.Type) reflect.Type {
	for fieldType != nil && fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	return fieldType
}

// walkSchema recursively compares a JSON-decoded value against a Go type,
// appending a SchemaMismatch for every object key with no matching struct field.
func walkSchema(raw any, fieldType reflect.Type, path, endpoint string, mismatches *[]SchemaMismatch) {
	if raw == nil || fieldType == nil {
		return
	}

	fieldType = indirectType(fieldType)

	switch fieldType.Kind() { //nolint:exhaustive // only structs, slices/arrays and maps carry nested keys worth checking
	case reflect.Struct:
		walkSchemaStruct(raw, fieldType, path, endpoint, mismatches)
	case reflect.Slice, reflect.Array:
		walkSchemaSlice(raw, fieldType.Elem(), path, endpoint, mismatches)
	case reflect.Map:
		walkSchemaMap(raw, fieldType.Elem(), path, endpoint, mismatches)
	default:
		// scalar, interface or other leaf kinds: nothing further to validate
	}
}

// walkSchemaStruct compares a JSON object against the fields declared on a
// Go struct type, recursing into every known field's value.
func walkSchemaStruct(raw any, structType reflect.Type, path, endpoint string, mismatches *[]SchemaMismatch) {
	object, ok := raw.(map[string]any)
	if !ok {
		return
	}

	knownFields := schemaFields(structType)

	for key, value := range object {
		knownFieldType, known := knownFields[key]
		if !known {
			*mismatches = append(*mismatches, SchemaMismatch{
				Endpoint: endpoint,
				GoType:   structType.String(),
				Path:     joinSchemaPath(path, key),
			})

			continue
		}

		walkSchema(value, knownFieldType, joinSchemaPath(path, key), endpoint, mismatches)
	}
}

// walkSchemaSlice recurses into every element of a JSON array against a Go
// slice/array element type.
func walkSchemaSlice(raw any, elemType reflect.Type, path, endpoint string, mismatches *[]SchemaMismatch) {
	array, ok := raw.([]any)
	if !ok {
		return
	}

	for _, item := range array {
		walkSchema(item, elemType, path+"[]", endpoint, mismatches)
	}
}

// walkSchemaMap recurses into every value of a JSON object against a Go map
// element type. Map keys themselves are never flagged as mismatches since
// maps intentionally accept arbitrary keys.
func walkSchemaMap(raw any, elemType reflect.Type, path, endpoint string, mismatches *[]SchemaMismatch) {
	object, ok := raw.(map[string]any)
	if !ok {
		return
	}

	for key, value := range object {
		walkSchema(value, elemType, joinSchemaPath(path, key), endpoint, mismatches)
	}
}

// schemaFields returns the JSON key to Go field type mapping declared on
// struct type structType. Anonymous (embedded) struct fields are flattened
// into the parent's field set, matching how encoding/json treats them.
func schemaFields(structType reflect.Type) map[string]reflect.Type {
	knownFields := make(map[string]reflect.Type, structType.NumField())

	for field := range structType.Fields() {
		name, embedded := schemaFieldName(field)

		// Unexported anonymous struct fields still promote their exported
		// inner fields (encoding/json does the same), so the PkgPath check
		// only applies to ordinary, non-embedded fields.
		switch {
		case embedded:
			mergeEmbeddedSchemaFields(knownFields, field.Type)
		case field.PkgPath != "":
			continue
		case name != "":
			knownFields[name] = field.Type
		}
	}

	return knownFields
}

// schemaFieldName returns the JSON key for field, and whether field is an
// embedded struct whose own fields should be flattened into the parent
// instead of being registered under a key of their own.
func schemaFieldName(field reflect.StructField) (name string, embedded bool) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false
	}

	name, _, _ = strings.Cut(tag, ",")

	if name == "" {
		if field.Anonymous {
			return "", true
		}

		name = field.Name
	}

	return name, false
}

// mergeEmbeddedSchemaFields flattens the fields of an embedded struct type
// into knownFields, mirroring encoding/json's handling of anonymous fields.
func mergeEmbeddedSchemaFields(knownFields map[string]reflect.Type, embeddedType reflect.Type) {
	embeddedType = indirectType(embeddedType)
	if embeddedType == nil || embeddedType.Kind() != reflect.Struct {
		return
	}

	maps.Copy(knownFields, schemaFields(embeddedType))
}

// joinSchemaPath appends key to path, separating existing segments with a dot.
func joinSchemaPath(path, key string) string {
	if path == "" {
		return key
	}

	return path + "." + key
}
