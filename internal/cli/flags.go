package cli

import (
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// BindFlags binds Cobra flags to a struct's fields based on reflection.
// It reads `url` struct tags to derive flag names and detects required fields
// (those without "omitempty" in their url tag).
// Embedded structs with `url:",inline"` or anonymous fields are recursively bound.
func BindFlags(cmd *cobra.Command, opt any) {
	bindFlagsRecursive(cmd.Flags(), reflect.ValueOf(opt).Elem(), "")
}

// bindFlagsRecursive walks the struct fields and creates appropriate pflags.
func bindFlagsRecursive(flags *pflag.FlagSet, structVal reflect.Value, prefix string) {
	structType := structVal.Type()

	for fieldIdx := range structType.NumField() {
		field := structType.Field(fieldIdx)
		fieldVal := structVal.Field(fieldIdx)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		urlTag := field.Tag.Get("url")

		// Handle embedded structs (e.g., PaginationArgs with url:",inline" or anonymous).
		if field.Anonymous || urlTag == ",inline" {
			if fieldVal.Kind() == reflect.Struct {
				bindFlagsRecursive(flags, fieldVal, prefix)

				continue
			}
		}

		// Skip fields with no url tag.
		if urlTag == "" {
			continue
		}

		flagName, required := parseFlagMeta(field)
		if prefix != "" {
			flagName = prefix + "-" + flagName
		}

		bindField(flags, fieldVal, field, flagName, required)
	}
}

// parseFlagMeta extracts flag name and required status from a struct field.
// Flag name is derived from the Go field name (PascalCase→kebab-case).
// A field is required if its url tag does not contain "omitempty".
func parseFlagMeta(field reflect.StructField) (string, bool) {
	flagName := pascalToKebab(field.Name)

	urlTag := field.Tag.Get("url")
	parts := strings.Split(urlTag, ",")
	required := !slices.Contains(parts[1:], "omitempty")

	return flagName, required
}

// bindField registers a pflag for the given struct field.
//
//nolint:exhaustive // only handling field types that appear in option structs
func bindField(flags *pflag.FlagSet, fieldVal reflect.Value, field reflect.StructField, flagName string, required bool) {
	description := buildFlagDescription(field, required)

	switch field.Type.Kind() {
	case reflect.String:
		stringPtr, _ := fieldVal.Addr().Interface().(*string)
		flags.StringVar(stringPtr, flagName, "", description)

	case reflect.Bool:
		boolPtr, _ := fieldVal.Addr().Interface().(*bool)
		flags.BoolVar(boolPtr, flagName, false, description)

	case reflect.Int64:
		int64Ptr, _ := fieldVal.Addr().Interface().(*int64)
		flags.Int64Var(int64Ptr, flagName, 0, description)

	case reflect.Slice:
		bindSliceField(flags, fieldVal, field, flagName, description)

	case reflect.Ptr:
		bindPointerField(flags, fieldVal, flagName, description)

	case reflect.Map:
		bindMapField(flags, fieldVal, field, flagName, description)
	}

	if required {
		_ = cobra.MarkFlagRequired(flags, flagName)
	}
}

// bindSliceField binds a []string field to a StringSlice flag.
func bindSliceField(flags *pflag.FlagSet, fieldVal reflect.Value, field reflect.StructField, flagName, description string) {
	if field.Type.Elem().Kind() == reflect.String {
		slicePtr, _ := fieldVal.Addr().Interface().(*[]string)
		flags.StringSliceVar(slicePtr, flagName, nil, description)
	}
}

// bindPointerField binds a *bool field to a TriStateBool custom flag.
func bindPointerField(flags *pflag.FlagSet, fieldVal reflect.Value, flagName, description string) {
	if fieldVal.Type().Elem().Kind() == reflect.Bool {
		target, _ := fieldVal.Addr().Interface().(**bool)
		flags.Var(NewTriStateBool(target), flagName, description)
	}
}

// bindMapField binds map fields using custom pflag.Value implementations.
//
//nolint:exhaustive // only handling map value types that appear in option structs
func bindMapField(flags *pflag.FlagSet, fieldVal reflect.Value, field reflect.StructField, flagName, description string) {
	keyType := field.Type.Key()
	elemType := field.Type.Elem()

	// Only handle maps with string keys.
	if keyType.Kind() != reflect.String {
		return
	}

	switch elemType.Kind() {
	case reflect.String:
		// map[string]string or named types like SemicolonSeparatedMap.
		// We allocate a plain map[string]string for flag binding, then copy into the field.
		mapPtr := newStringMapPtr(fieldVal)
		flags.Var(NewStringMapValue(mapPtr), flagName, description)

	case reflect.Interface:
		// map[string]any or named types like JSONEncodedMap.
		mapPtr := newAnyMapPtr(fieldVal)
		flags.Var(NewJSONMapValue(mapPtr), flagName, description)
	}
}

// newStringMapPtr returns a *map[string]string that, when written to, also updates
// the original reflect.Value (which may be a named map type like SemicolonSeparatedMap).
func newStringMapPtr(fieldVal reflect.Value) *map[string]string {
	plainMapType := reflect.TypeFor[map[string]string]()
	emptyMap := make(map[string]string)
	fieldVal.Set(reflect.ValueOf(emptyMap).Convert(fieldVal.Type()))

	//nolint:forcetypeassert // type is guaranteed by caller check on elemType.Kind()
	plain := fieldVal.Convert(plainMapType).Interface().(map[string]string)

	return &plain
}

// newAnyMapPtr returns a *map[string]any that, when written to, also updates
// the original reflect.Value (which may be a named map type like JSONEncodedMap).
func newAnyMapPtr(fieldVal reflect.Value) *map[string]any {
	plainMapType := reflect.TypeFor[map[string]any]()
	emptyMap := make(map[string]any)
	fieldVal.Set(reflect.ValueOf(emptyMap).Convert(fieldVal.Type()))

	//nolint:forcetypeassert // type is guaranteed by caller check on elemType.Kind()
	plain := fieldVal.Convert(plainMapType).Interface().(map[string]any)

	return &plain
}

// buildFlagDescription generates a flag description from the struct field comment and url tag.
func buildFlagDescription(field reflect.StructField, required bool) string {
	urlTag := field.Tag.Get("url")
	parts := strings.Split(urlTag, ",")
	apiParam := parts[0]

	desc := "API parameter: " + apiParam

	if required {
		desc += " (required)"
	}

	return desc
}

// pascalToKebab converts a PascalCase string to kebab-case.
// e.g., "ImpactSeverities" → "impact-severities", "PciDss32" → "pci-dss32".
// Special cases: "L10N" → "l10n".
//

func pascalToKebab(input string) string {
	if input == "" {
		return input
	}

	var result strings.Builder

	runes := []rune(input)
	for idx, char := range runes {
		if unicode.IsUpper(char) {
			if shouldInsertDash(runes, idx) {
				result.WriteRune('-')
			}

			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// shouldInsertDash decides if a dash should be inserted before an uppercase rune at the given index.
func shouldInsertDash(runes []rune, idx int) bool {
	if idx == 0 {
		return false
	}

	prev := runes[idx-1]

	// After a lowercase letter: aB → a-b
	if unicode.IsLower(prev) {
		return true
	}

	// After an uppercase letter and before a lowercase letter: ABc → a-bc
	// NOT after a digit: 10N should not become 10-n
	if unicode.IsUpper(prev) && idx+1 < len(runes) && unicode.IsLower(runes[idx+1]) {
		return true
	}

	return false
}
