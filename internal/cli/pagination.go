package cli

import (
	"reflect"
	"strings"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

const (
	// paginationPageField is the field name for the page number.
	paginationPageField = "Page"
	// paginationPageSizeField is the field name for the page size.
	paginationPageSizeField = "PageSize"
)

// hasPagination checks if an option struct embeds PaginationArgs.
func hasPagination(optType reflect.Type) bool {
	for field := range optType.Fields() {
		if field.Anonymous && field.Type.Name() == "PaginationArgs" {
			return true
		}
	}

	return false
}

// responseHasPaging checks if a response struct has a Paging field.
func responseHasPaging(responseType reflect.Type) bool {
	for responseType.Kind() == reflect.Pointer {
		responseType = responseType.Elem()
	}

	if responseType.Kind() != reflect.Struct {
		return false
	}

	_, found := responseType.FieldByName("Paging")

	return found
}

// facetsFieldName is the name of the facets slice found on search responses.
// Facets are aggregated metadata, never the paginated items themselves.
const facetsFieldName = "Facets"

// findSliceField returns the name of the exported []struct field that holds the paginated items
// of a response type (e.g., Issues, Components, Projects).
//
// Many responses carry several []struct fields alongside the paginated ones (e.g., IssuesSearch has
// Components, Facets, Issues, Rules and Users), and Paging.Total only describes the paginated one.
// The field is chosen as follows:
//  1. the only candidate, if there is exactly one;
//  2. the candidate whose name prefixes the response type name (IssuesSearch -> Issues);
//  3. the first candidate that is not Facets;
//  4. the first candidate.
func findSliceField(responseType reflect.Type) (string, bool) {
	for responseType.Kind() == reflect.Pointer {
		responseType = responseType.Elem()
	}

	if responseType.Kind() != reflect.Struct {
		return "", false
	}

	candidates := structSliceFields(responseType)

	switch len(candidates) {
	case 0:
		return "", false
	case 1:
		return candidates[0], true
	}

	typeName := responseType.Name()
	for _, name := range candidates {
		if strings.HasPrefix(typeName, name) {
			return name, true
		}
	}

	for _, name := range candidates {
		if name != facetsFieldName {
			return name, true
		}
	}

	return candidates[0], true
}

// structSliceFields returns the names of all exported []struct (or []*struct) fields, in declaration order.
func structSliceFields(structType reflect.Type) []string {
	var names []string

	for field := range structType.Fields() {
		if !field.IsExported() || field.Type.Kind() != reflect.Slice {
			continue
		}

		elemType := field.Type.Elem()
		for elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}

		if elemType.Kind() == reflect.Struct {
			names = append(names, field.Name)
		}
	}

	return names
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
		// No slice field to paginate - just call once.
		result, resp, err := InvokeMethod(service, methodName, opt, pattern, true)
		CloseBody(resp)

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

		result, resp, err := InvokeMethod(service, methodName, opt, pattern, true)
		CloseBody(resp)

		if err != nil {
			return nil, err
		}

		if result == nil {
			break
		}

		resultVal := reflect.ValueOf(result)
		for resultVal.Kind() == reflect.Pointer {
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

		if isLastPage(resultVal, itemsField.Len(), pageNum) {
			break
		}

		pageNum++
	}

	// Set the accumulated items on the first result and return it.
	setPaginatedResult(firstResult, allItems, sliceFieldName)

	return firstResult, nil
}

// isLastPage reports whether no further page should be requested after the given page.
// It stops on an empty page or once the pages fetched so far cover Paging.Total, so a response
// whose item count disagrees with Paging.Total can never cause an endless loop.
func isLastPage(resultVal reflect.Value, pageItems int, pageNum int64) bool {
	if pageItems == 0 {
		return true
	}

	pagingField := resultVal.FieldByName("Paging")
	if !pagingField.IsValid() {
		return true
	}

	pageSize := pagingField.FieldByName("PageSize").Int()
	if pageSize <= 0 {
		pageSize = int64(sonar.MaxPageSize)
	}

	return pageNum*pageSize >= pagingField.FieldByName("Total").Int()
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
	for resultVal.Kind() == reflect.Pointer {
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
