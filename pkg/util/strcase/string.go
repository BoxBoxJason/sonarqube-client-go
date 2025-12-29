package strcase

import (
	"strings"
)

// toCamelInitCase converts a string to CamelCase.
//
//nolint:cyclop
func toCamelInitCase(str string, initCase bool) string {
	str = addWordBoundariesToNumbers(str)
	str = strings.Trim(str, " ")
	result := ""
	capNext := initCase

	var nSb13 strings.Builder

	for _, char := range str {
		if char >= 'A' && char <= 'Z' {
			nSb13.WriteRune(char)
		}

		if char >= '0' && char <= '9' {
			nSb13.WriteRune(char)
		}

		if char >= 'a' && char <= 'z' {
			if capNext {
				nSb13.WriteString(strings.ToUpper(string(char)))
			} else {
				nSb13.WriteRune(char)
			}
		}

		if char == '_' || char == ' ' || char == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}

	result += nSb13.String()

	return result
}

// ToCamel converts a string to CamelCase.
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel converts a string to lowerCamelCase.
func ToLowerCamel(str string) string {
	if str == "" {
		return str
	}

	if r := rune(str[0]); r >= 'A' && r <= 'Z' {
		str = strings.ToLower(string(r)) + str[1:]
	}

	return toCamelInitCase(str, false)
}
