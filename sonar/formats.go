package sonargo

import (
	"net/url"
	"strings"
)

// ListToSeparatedString converts a list of strings into a single string separated by the given separator.
func ListToSeparatedString(list []string, separator string) string {
	result := ""
	for i, item := range list {
		if i > 0 {
			result += separator
		}
		result += item
	}
	return result
}

// SeparatedStringToList converts a separated string into a list of strings using the given separator.
func SeparatedStringToList(s string, separator string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, separator)
}

// MapToSeparatedString converts a map of strings into a single string with entries separated by entrySeparator
// and key-value pairs separated by keyValueSeparator.
func MapToSeparatedString(m map[string]string, entrySeparator string, keyValueSeparator string) string {
	result := ""
	i := 0
	for k, v := range m {
		if i > 0 {
			result += entrySeparator
		}
		result += k + keyValueSeparator + v
		i++
	}
	return result
}

// SeparatedStringToMap converts a separated string into a map of strings using the given entrySeparator and keyValueSeparator.
func SeparatedStringToMap(s string, entrySeparator string, keyValueSeparator string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	entries := strings.Split(s, entrySeparator)
	for _, entry := range entries {
		parts := strings.SplitN(entry, keyValueSeparator, 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

// CommaSeparatedSlice is a custom type for URL encoding string slices as comma-separated values.
type CommaSeparatedSlice []string

// EncodeValues implements the query.Encoder interface for custom URL encoding.
func (s CommaSeparatedSlice) EncodeValues(key string, v *url.Values) error {
	if len(s) > 0 {
		v.Set(key, ListToSeparatedString(s, ","))
	}
	return nil
}

// SemicolonSeparatedMap is a custom type for URL encoding maps as semicolon-separated key=value pairs.
type SemicolonSeparatedMap map[string]string

// EncodeValues implements the query.Encoder interface for custom URL encoding.
func (m SemicolonSeparatedMap) EncodeValues(key string, v *url.Values) error {
	if len(m) > 0 {
		v.Set(key, MapToSeparatedString(m, ";", "="))
	}
	return nil
}

// EncodeSliceToCommaSeparated converts a string slice to a comma-separated URL value.
// This is a helper for encoding []string fields that need to be sent as comma-separated values.
func EncodeSliceToCommaSeparated(key string, values []string) string {
	if len(values) == 0 {
		return ""
	}

	return key + "=" + url.QueryEscape(ListToSeparatedString(values, ","))
}

// EncodeMapToSeparated converts a map[string]string to a URL value with custom separators.
// This is a helper for encoding map fields that need to be sent in a specific format.
// For example, impacts=MAINTAINABILITY=HIGH;SECURITY=LOW.
func EncodeMapToSeparated(key string, m map[string]string, entrySep, kvSep string) string {
	if len(m) == 0 {
		return ""
	}
	return key + "=" + url.QueryEscape(MapToSeparatedString(m, entrySep, kvSep))
}
