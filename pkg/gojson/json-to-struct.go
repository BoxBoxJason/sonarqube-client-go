// Package gojson generates go struct defintions from JSON documents
//
// # Reads from stdin and prints to stdout
package gojson

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

// ForceFloats forces all numbers to be floats.
var ForceFloats bool //nolint:gochecknoglobals

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{ //nolint:gochecknoglobals
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"NTP":   true,
	"DB":    true,
}

var intToWordMap = []string{ //nolint:gochecknoglobals
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// Parser is a function that parses input into an interface.
type Parser func(io.Reader) (any, error)

// ParseJSON parses JSON input.
func ParseJSON(input io.Reader) (any, error) {
	var (
		result any
		json   = jsoniter.ConfigCompatibleWithStandardLibrary
	)

	byts, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	err = json.Unmarshal(byts, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return result, nil
}

// ParseYAML parses YAML input.
func ParseYAML(input io.Reader) (any, error) {
	var result any

	content, err := readFile(input)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	return result, nil
}

//nolint:unparam
func readFile(input io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	_, err := io.Copy(buf, input)
	if err != nil {
		return []byte{}, nil //nolint:nilerr // This seems to be intended behavior in original code
	}

	return buf.Bytes(), nil
}

// Generate a struct definition given a JSON string representation of an object and a name structName.
//
//nolint:cyclop
func Generate(input io.Reader, parser Parser, structName string, tags []string, subStruct bool, convertFloats bool) ([]byte, error) {
	var subStructMap map[string]string
	if subStruct {
		subStructMap = make(map[string]string)
	}

	var result map[string]any

	iresult, err := parser(input)
	if err != nil {
		return nil, err
	}

	switch iresult := iresult.(type) {
	case map[any]any:
		result = convertKeysToStrings(iresult)
	case map[string]any:
		result = iresult
	case []any:
		src := fmt.Sprintf("type %s %s\n",
			structName,
			typeForValue(iresult, structName, tags, subStructMap, convertFloats))

		formatted, fmtErr := format.Source([]byte(src))
		if fmtErr != nil {
			return nil, fmt.Errorf("error formatting: %w, was formatting\n%s", fmtErr, src)
		}

		return formatted, nil
	default:
		return nil, fmt.Errorf("unexpected type: %T", iresult)
	}

	src := fmt.Sprintf("type %s %s}",
		structName,
		generateTypes(result, structName, tags, 0, subStructMap, convertFloats))

	keys := make([]string, 0, len(subStructMap))
	for key := range subStructMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, k := range keys {
		src = fmt.Sprintf("%v\n\ntype %v %v", src, subStructMap[k], k)
	}

	formatted, err := format.Source([]byte(src))
	if err != nil {
		return nil, fmt.Errorf("error formatting: %w, was formatting\n%s", err, src)
	}

	return formatted, nil
}

func convertKeysToStrings(obj map[any]any) map[string]any {
	res := make(map[string]any)

	for k, v := range obj {
		res[fmt.Sprintf("%v", k)] = v
	}

	return res
}

// generateTypes Generate go struct entries for a map[string]any structure
//
//nolint:gocognit,cyclop,funlen
func generateTypes(obj map[string]any, structName string, tags []string, depth int, subStructMap map[string]string, convertFloats bool) string {
	var builder strings.Builder

	builder.WriteString("struct {")

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		value := obj[key]
		valueType := typeForValue(value, structName, tags, subStructMap, convertFloats)

		// If a nested value, recurse
		switch value := value.(type) {
		case []any:
			if len(value) > 0 { //nolint:nestif
				sub := ""
				if v, ok := value[0].(map[any]any); ok {
					sub = generateTypes(convertKeysToStrings(v), structName, tags, depth+1, subStructMap, convertFloats) + "}"
				} else if v, ok := value[0].(map[string]any); ok {
					sub = generateTypes(v, structName, tags, depth+1, subStructMap, convertFloats) + "}"
				}

				if sub != "" {
					subName := sub

					if subStructMap != nil {
						if val, ok := subStructMap[sub]; ok {
							subName = val
						} else {
							subName = fmt.Sprintf("%v_sub%v", structName, len(subStructMap)+1)

							subStructMap[sub] = subName
						}
					}

					valueType = "[]" + subName
				}
			}
		case map[any]any:
			sub := generateTypes(convertKeysToStrings(value), structName, tags, depth+1, subStructMap, convertFloats) + "}"
			subName := sub

			if subStructMap != nil {
				if val, ok := subStructMap[sub]; ok {
					subName = val
				} else {
					subName = fmt.Sprintf("%v_sub%v", structName, len(subStructMap)+1)

					subStructMap[sub] = subName
				}
			}

			valueType = subName
		case map[string]any:
			sub := generateTypes(value, structName, tags, depth+1, subStructMap, convertFloats) + "}"
			subName := sub

			if subStructMap != nil {
				if val, ok := subStructMap[sub]; ok {
					subName = val
				} else {
					subName = fmt.Sprintf("%v_sub%v", structName, len(subStructMap)+1)

					subStructMap[sub] = subName
				}
			}

			valueType = subName
		}

		fieldName := FmtFieldName(key)
		tagList := make([]string, len(tags))

		for i, t := range tags {
			tagList[i] = fmt.Sprintf("%s:\"%s,omitempty\"", t, key)
		}

		builder.WriteString(fmt.Sprintf("\n%s %s `%s`",
			fieldName,
			valueType,
			strings.Join(tagList, " ")))
	}

	return builder.String()
}

// FmtFieldName formats a string as a struct key
//
// Example:
//
// FmtFieldName("foo_id")/
//
// Output: FooID.
func FmtFieldName(name string) string {
	runes := []rune(name)
	for len(runes) > 0 && !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
		runes = runes[1:]
	}

	if len(runes) == 0 {
		return "_"
	}

	name = stringifyFirstChar(string(runes))
	lintedName := lintFieldName(name)
	runes = []rune(lintedName)

	for idx, char := range runes {
		ok := unicode.IsLetter(char) || unicode.IsDigit(char)
		if idx == 0 {
			ok = unicode.IsLetter(char)
		}

		if !ok {
			runes[idx] = '_'
		}
	}

	name = string(runes)
	name = strings.Trim(name, "_")

	if len(name) == 0 {
		return "_"
	}

	return name
}

//nolint:gocognit,cyclop,funlen,gocritic
func lintFieldName(name string) string {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}

	allLower := true

	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false

			break
		}
	}

	if allLower {
		runes := []rune(name)
		if u := strings.ToUpper(name); commonInitialisms[u] {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}

		return string(runes)
	}

	allUpperWithUnderscore := true

	for _, r := range name {
		if !unicode.IsUpper(r) && r != '_' {
			allUpperWithUnderscore = false

			break
		}
	}

	if allUpperWithUnderscore {
		name = strings.ToLower(name)
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	wordStart, scanIndex := 0, 0 // index of start of word, scan

	for scanIndex+1 <= len(runes) {
		eow := false // whether we hit the end of a word

		if scanIndex+1 == len(runes) {
			eow = true
		} else if runes[scanIndex+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			underscoreCount := 1

			for scanIndex+underscoreCount+1 < len(runes) && runes[scanIndex+underscoreCount+1] == '_' {
				underscoreCount++
			}

			// Leave at most one underscore if the underscore is between two digits
			if scanIndex+underscoreCount+1 < len(runes) && unicode.IsDigit(runes[scanIndex]) && unicode.IsDigit(runes[scanIndex+underscoreCount+1]) {
				underscoreCount--
			}

			copy(runes[scanIndex+1:], runes[scanIndex+underscoreCount+1:])
			runes = runes[:len(runes)-underscoreCount]
		} else if unicode.IsLower(runes[scanIndex]) && !unicode.IsLower(runes[scanIndex+1]) {
			// lower->non-lower
			eow = true
		}

		scanIndex++

		if !eow {
			continue
		}

		// [wordStart,scanIndex) is a word.
		word := string(runes[wordStart:scanIndex])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			copy(runes[wordStart:], []rune(u))
		} else if strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[wordStart] = unicode.ToUpper(runes[wordStart])
		}

		wordStart = scanIndex
	}

	return string(runes)
}

// generate an appropriate struct type entry.
func typeForValue(value any, structName string, tags []string, subStructMap map[string]string, convertFloats bool) string {
	// Check if this is an array
	if objects, ok := value.([]any); ok {
		types := make(map[reflect.Type]bool, 0)
		for _, o := range objects {
			types[reflect.TypeOf(o)] = true
		}

		if len(types) == 1 {
			//nolint:forcetypeassert // We know it's []any
			return "[]" + typeForValue(mergeElements(objects).([]any)[0], structName, tags, subStructMap, convertFloats)
		}

		return "[]any"
	} else if object, ok := value.(map[any]any); ok {
		return generateTypes(convertKeysToStrings(object), structName, tags, 0, subStructMap, convertFloats) + "}"
	} else if object, ok := value.(map[string]any); ok {
		return generateTypes(object, structName, tags, 0, subStructMap, convertFloats) + "}"
	} else if reflect.TypeOf(value) == nil {
		return "any"
	}

	v := reflect.TypeOf(value).Name()
	if v == "float64" && convertFloats {
		v = disambiguateFloatInt(value)
	}

	return v
}

// All numbers will initially be read as float64
// If the number appears to be an integer value, use int instead.
func disambiguateFloatInt(value any) string {
	const epsilon = .0001

	vfloat, ok := value.(float64)
	if !ok {
		return reflect.TypeOf(value).Name()
	}

	if !ForceFloats && math.Abs(vfloat-math.Floor(vfloat+epsilon)) < epsilon {
		var tmp int64

		return reflect.TypeOf(tmp).Name() //nolint:modernize
	}

	return reflect.TypeOf(value).Name()
}

// convert first character ints to strings.
func stringifyFirstChar(str string) string {
	first := str[:1]

	i, err := strconv.ParseInt(first, 10, 8)
	if err != nil {
		return str
	}

	return intToWordMap[i] + "_" + str[1:]
}

func mergeElements(val any) any {
	switch val := val.(type) {
	default:
		return val
	case []any:
		l := len(val)
		if l == 0 {
			return val
		}

		for j := 1; j < l; j++ {
			val[0] = mergeObjects(val[0], val[j])
		}

		return val[0:1]
	}
}

//nolint:cyclop,gocritic
func mergeObjects(obj1, obj2 any) any {
	if obj1 == nil {
		return obj2
	}

	if obj2 == nil {
		return obj1
	}

	if reflect.TypeOf(obj1) != reflect.TypeOf(obj2) {
		return nil
	}

	switch val := obj1.(type) {
	default:
		return obj1
	case []any:
		if i2, ok := obj2.([]any); ok {
			i3 := append(val, i2...)

			return mergeElements(i3)
		}

		return mergeElements(val)
	case map[string]any:
		if i2, ok := obj2.(map[string]any); ok {
			for k, v := range i2 {
				if v2, ok := val[k]; ok {
					val[k] = mergeObjects(v2, v)
				} else {
					val[k] = v
				}
			}
		}

		return val
	case map[any]any:
		if i2, ok := obj2.(map[any]any); ok {
			for k, v := range i2 {
				if v2, ok := val[k]; ok {
					val[k] = mergeObjects(v2, v)
				} else {
					val[k] = v
				}
			}
		}

		return val
	}
}
