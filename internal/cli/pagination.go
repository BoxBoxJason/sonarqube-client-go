package cli

import (
	"reflect"

	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

const (
	// paginationPageField is the field name for the page number.
	paginationPageField = "Page"
	// paginationPageSizeField is the field name for the page size.
	paginationPageSizeField = "PageSize"
)

// hasPagination checks if an option struct embeds PaginationArgs.
func hasPagination(optType reflect.Type) bool {
	for i := range optType.NumField() {
		field := optType.Field(i)
		if field.Anonymous && field.Type.Name() == "PaginationArgs" {
			return true
		}
	}

	return false
}

// responseHasPaging checks if a response struct has a Paging field.
func responseHasPaging(responseType reflect.Type) bool {
	for responseType.Kind() == reflect.Ptr {
		responseType = responseType.Elem()
	}

	if responseType.Kind() != reflect.Struct {
		return false
	}

	_, found := responseType.FieldByName("Paging")

	return found
}

// findSliceField returns the name and type of the first exported []struct field in a response type.
// This is the field that contains the paginated items (e.g., Issues, Components, Projects).
func findSliceField(responseType reflect.Type) (string, bool) {
	for responseType.Kind() == reflect.Ptr {
		responseType = responseType.Elem()
	}

	for i := range responseType.NumField() {
		field := responseType.Field(i)
		if !field.IsExported() {
			continue
		}

		if field.Type.Kind() == reflect.Slice {
			elemType := field.Type.Elem()
			for elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			if elemType.Kind() == reflect.Struct {
				return field.Name, true
			}
		}
	}

	return "", false
}

// PaginateAll calls a paginated service method repeatedly until all results are collected.
// It sets Page=1 and PageSize=MaxPageSize, then increments page until Total is reached.
// The results are merged by concatenating the slice field across pages.
//
//nolint:funlen // unavoidable length for multi-page merging with reflection
func PaginateAll(
	service reflect.Value,
	methodName string,
	opt reflect.Value,
	pattern MethodReturnPattern,
	responseType reflect.Type,
) (any, error) {
	sliceFieldName, hasSlice := findSliceField(responseType)
	if !hasSlice {
		// No slice field to paginate — just call once.
		result, _, err := InvokeMethod(service, methodName, opt, pattern, true) //nolint:bodyclose // caller manages lifecycle

		return result, err
	}

	// Set initial pagination values.
	optElem := opt.Elem()
	setPageField(optElem, paginationPageField, 1)
	setPageField(optElem, paginationPageSizeField, int64(sonar.MaxPageSize))

	var (
		allItems    reflect.Value
		firstResult any
		pageNum     int64 = 1
	)

	for {
		setPageField(optElem, paginationPageField, pageNum)

		result, _, err := InvokeMethod(service, methodName, opt, pattern, true) //nolint:bodyclose // caller manages lifecycle
		if err != nil {
			return nil, err
		}

		if result == nil {
			break
		}

		resultVal := reflect.ValueOf(result)
		for resultVal.Kind() == reflect.Ptr {
			resultVal = resultVal.Elem()
		}

		// Get the slice field from this page's result.
		itemsField := resultVal.FieldByName(sliceFieldName)
		if !itemsField.IsValid() {
			break
		}

		if pageNum == 1 {
			firstResult = result
			allItems = reflect.MakeSlice(itemsField.Type(), 0, itemsField.Len())
		}

		allItems = reflect.AppendSlice(allItems, itemsField)

		// Check if we've fetched all items.
		pagingField := resultVal.FieldByName("Paging")
		if !pagingField.IsValid() {
			break
		}

		total := pagingField.FieldByName("Total").Int()

		if int64(allItems.Len()) >= total {
			break
		}

		pageNum++
	}

	// Set the accumulated items on the first result and return it.
	setPaginatedResult(firstResult, allItems, sliceFieldName)

	return firstResult, nil
}

// setPageField sets a pagination field (Page or PageSize) on an option struct
// by navigating into the embedded PaginationArgs.
func setPageField(optElem reflect.Value, fieldName string, value int64) {
	// Try direct field first.
	field := optElem.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		field.SetInt(value)
	}
}

// setPaginatedResult sets the accumulated items on the first result and updates paging info.
func setPaginatedResult(firstResult any, allItems reflect.Value, sliceFieldName string) {
	if firstResult == nil || !allItems.IsValid() {
		return
	}

	resultVal := reflect.ValueOf(firstResult)
	for resultVal.Kind() == reflect.Ptr {
		resultVal = resultVal.Elem()
	}

	resultVal.FieldByName(sliceFieldName).Set(allItems)

	// Update paging info to reflect total.
	pagingField := resultVal.FieldByName("Paging")
	if pagingField.IsValid() {
		pagingField.FieldByName("PageIndex").SetInt(1)
		pagingField.FieldByName("PageSize").SetInt(int64(allItems.Len()))
	}
}
