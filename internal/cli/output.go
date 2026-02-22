package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// OutputFormat represents the supported output formats for CLI responses.
type OutputFormat string

const (
	// OutputJSON formats output as indented JSON.
	OutputJSON OutputFormat = "json"
	// OutputTable formats output as an ASCII table.
	OutputTable OutputFormat = "table"
	// OutputYAML formats output as YAML.
	OutputYAML OutputFormat = "yaml"

	// tablePadding is the padding added to each column in table output.
	tablePadding = 2
)

// FormatOutput writes the given value to the writer in the specified format.
// It handles nil values, raw bytes, raw strings, and struct/slice types.
func FormatOutput(writer io.Writer, val any, format OutputFormat) error {
	if val == nil {
		return nil
	}

	written, err := writeRawValue(writer, val)
	if written || err != nil {
		return err
	}

	switch format {
	case OutputJSON:
		return formatJSON(writer, val)
	case OutputYAML:
		return formatYAML(writer, val)
	case OutputTable:
		return formatTable(writer, val)
	default:
		return formatJSON(writer, val)
	}
}

// writeRawValue handles raw byte and raw string output directly.
// Returns true if the value was handled, false if it should be formatted normally.
func writeRawValue(writer io.Writer, val any) (bool, error) {
	// Handle raw byte output (e.g., from Batch.File).
	if byteSlice, ok := val.([]byte); ok {
		_, err := writer.Write(byteSlice)
		if err != nil {
			return true, fmt.Errorf("failed to write raw bytes: %w", err)
		}

		return true, nil
	}

	// Handle raw string pointer output (e.g., from Batch.Index).
	if strPtr, ok := val.(*string); ok && strPtr != nil {
		_, err := fmt.Fprintln(writer, *strPtr)
		if err != nil {
			return true, fmt.Errorf("failed to write string: %w", err)
		}

		return true, nil
	}

	return false, nil
}

// formatJSON outputs the value as indented JSON.
func formatJSON(writer io.Writer, data any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// formatYAML outputs the value as YAML.
func formatYAML(writer io.Writer, data any) error {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	_, err = writer.Write(yamlBytes)
	if err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}

	return nil
}

// formatTable outputs the value as an ASCII table.
// For slices of structs, each struct field becomes a column.
// For single structs, output as a key-value table.
func formatTable(writer io.Writer, data any) error {
	rval := reflect.ValueOf(data)

	// Dereference pointers.
	for rval.Kind() == reflect.Ptr {
		if rval.IsNil() {
			return nil
		}

		rval = rval.Elem()
	}

	//nolint:exhaustive // only struct and slice are renderable as tables
	switch rval.Kind() {
	case reflect.Slice:
		return formatSliceTable(writer, rval)
	case reflect.Struct:
		return formatStructTable(writer, rval)
	default:
		// Fallback to JSON for unsupported types.
		return formatJSON(writer, data)
	}
}

// formatSliceTable formats a slice of structs as a table with columns.
func formatSliceTable(writer io.Writer, sliceVal reflect.Value) error {
	if sliceVal.Len() == 0 {
		_, err := fmt.Fprintln(writer, "(no results)")
		if err != nil {
			return fmt.Errorf("failed to write empty result: %w", err)
		}

		return nil
	}

	// Check if elements are structs.
	elemType := sliceVal.Type().Elem()
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return formatJSON(writer, sliceVal.Interface())
	}

	// Extract headers from struct fields.
	headers := extractHeaders(elemType)

	// Extract rows.
	rows := make([][]string, sliceVal.Len())
	for rowIdx := range sliceVal.Len() {
		elem := sliceVal.Index(rowIdx)

		for elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		rows[rowIdx] = extractRow(elem, elemType)
	}

	renderTable(writer, headers, rows)

	return nil
}

// formatStructTable formats a single struct as a key-value table.
func formatStructTable(writer io.Writer, structVal reflect.Value) error {
	structType := structVal.Type()
	headers := []string{"FIELD", "VALUE"}
	rows := make([][]string, 0, structType.NumField())

	for fieldIdx := range structType.NumField() {
		field := structType.Field(fieldIdx)
		if !field.IsExported() {
			continue
		}

		fieldVal := structVal.Field(fieldIdx)

		// Skip zero values.
		if fieldVal.IsZero() {
			continue
		}

		// For nested structs/slices, format as JSON inline.
		var valueStr string

		//nolint:exhaustive // only struct/slice/map need JSON inline formatting
		switch fieldVal.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Map:
			data, err := json.Marshal(fieldVal.Interface())
			if err != nil {
				valueStr = fmt.Sprintf("%v", fieldVal.Interface())
			} else {
				valueStr = string(data)
			}
		default:
			valueStr = fmt.Sprintf("%v", fieldVal.Interface())
		}

		rows = append(rows, []string{field.Name, valueStr})
	}

	renderTable(writer, headers, rows)

	return nil
}

// extractHeaders returns the column headers from a struct type.
// Uses the JSON tag name if available, otherwise the field name.
func extractHeaders(structType reflect.Type) []string {
	headers := make([]string, 0, structType.NumField())

	for fieldIdx := range structType.NumField() {
		field := structType.Field(fieldIdx)
		if !field.IsExported() {
			continue
		}

		name := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}

		headers = append(headers, strings.ToUpper(name))
	}

	return headers
}

// extractRow returns the column values from a struct value.
func extractRow(structVal reflect.Value, structType reflect.Type) []string {
	row := make([]string, 0, structType.NumField())

	for fieldIdx := range structType.NumField() {
		field := structType.Field(fieldIdx)
		if !field.IsExported() {
			continue
		}

		fieldVal := structVal.Field(fieldIdx)

		var cellStr string

		//nolint:exhaustive // only struct/slice/map need JSON inline formatting
		switch fieldVal.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Map:
			if fieldVal.IsZero() {
				cellStr = ""
			} else {
				data, err := json.Marshal(fieldVal.Interface())
				if err != nil {
					cellStr = fmt.Sprintf("%v", fieldVal.Interface())
				} else {
					cellStr = string(data)
				}
			}
		default:
			if fieldVal.IsZero() {
				cellStr = ""
			} else {
				cellStr = fmt.Sprintf("%v", fieldVal.Interface())
			}
		}

		row = append(row, cellStr)
	}

	return row
}

// renderTable renders headers and rows as a simple ASCII table.
func renderTable(writer io.Writer, headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths.
	widths := make([]int, len(headers))
	for colIdx, header := range headers {
		widths[colIdx] = len(header)
	}

	for _, row := range rows {
		for colIdx, cell := range row {
			if colIdx < len(widths) && len(cell) > widths[colIdx] {
				widths[colIdx] = len(cell)
			}
		}
	}

	// Print header row.
	printTableRow(writer, headers, widths)

	// Print separator.
	sep := make([]string, len(widths))
	for colIdx, width := range widths {
		sep[colIdx] = strings.Repeat("-", width+tablePadding)
	}

	_, _ = fmt.Fprintln(writer, strings.Join(sep, "+"))

	// Print data rows.
	for _, row := range rows {
		printTableRow(writer, row, widths)
	}
}

// printTableRow prints a single table row with proper column widths.
func printTableRow(writer io.Writer, cells []string, widths []int) {
	parts := make([]string, len(widths))

	for colIdx := range widths {
		cell := ""
		if colIdx < len(cells) {
			cell = cells[colIdx]
		}

		parts[colIdx] = fmt.Sprintf("%-*s", widths[colIdx]+tablePadding, cell)
	}

	_, _ = fmt.Fprintln(writer, strings.Join(parts, "| "))
}
