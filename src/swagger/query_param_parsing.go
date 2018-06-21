package swagger

import (
	"fmt"
	"net/url"
	"strconv"
)

func readLastQueryParam(u *url.URL) (bool, error) {
	vars := u.Query()
	if err := checkUnexpectedQueryParams(&vars, "last"); err != nil {
		return false, err
	}

	last, err := readQueryParam(&vars, "last", false, true,
		func(value string) (interface{}, error) { return strconv.ParseBool(value) })
	return last.(bool), err
}

func readSizeAndFromQueryParams(u *url.URL, defaultSize int) (int, string, error) {
	vars := u.Query()
	if err := checkUnexpectedQueryParams(&vars, "size", "from"); err != nil {
		return defaultSize, "", err
	}

	size, err := readQueryParam(&vars, "size", defaultSize, nil,
		func(value string) (interface{}, error) { return strconv.Atoi(value) })
	if err != nil {
		return defaultSize, "", err
	}

	from, err := readQueryParam(&vars, "from", "", nil,
		func(value string) (interface{}, error) { return value, nil })
	return size.(int), from.(string), err
}

func checkUnexpectedQueryParams(vars *url.Values, expectedParams ...string) error {
	expectedParamsMap := make(map[string]struct{})
	for _, param := range expectedParams {
		expectedParamsMap[param] = struct{}{}
	}

	for key := range *vars {
		if _, found := expectedParamsMap[key]; !found {
			return fmt.Errorf("Unexpected query param '%v'", key)
		}
	}

	return nil
}

func readQueryParam(vars *url.Values, paramName string, defaultResult interface{}, emptyResult interface{}, parseFunc func(string) (interface{}, error)) (interface{}, error) {
	value, found := (*vars)[paramName]
	if !found {
		return defaultResult, nil
	}

	if len(value) > 1 {
		return defaultResult, fmt.Errorf("Multiple values specified for query param '%v'", paramName)
	}

	if value[0] == "" {
		if emptyResult != nil {
			return emptyResult, nil
		}
		return defaultResult, fmt.Errorf("The query param '%v' has no value specified", paramName)
	}

	result, err := parseFunc(value[0])
	if err != nil {
		return defaultResult, fmt.Errorf("The query param '%v' should be a(n) %T value, got '%v'", paramName, defaultResult, value[0])
	}

	return result, nil
}
