package swagger

import (
	"net/url"
	"testing"
)

type testValues struct {
	name          string
	url           string
	errorExpected bool
}

func (values *testValues) getTestFunc(funcToTest func(*url.URL) ([]interface{}, error), validate func([]interface{}, *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		testURL, err := url.Parse(values.url)
		if err != nil {
			t.Fatalf("Invalid test URL value '%v' - %v", values.url, err.Error())
		}

		results, err := funcToTest(testURL)
		if err != nil {
			if values.errorExpected {
				return
			}

			t.Fatalf("Unexpected error: %v", err.Error())
		}

		if values.errorExpected {
			t.Fatalf("Expected an error, but got none")
		}

		validate(results, t)
	}
}

type readLastQueryParamTestValues struct {
	test         testValues
	expectedLast bool
}

func (values *readLastQueryParamTestValues) getTestFunc() func(*testing.T) {
	return values.test.getTestFunc(
		func(testURL *url.URL) ([]interface{}, error) {
			last, err := readLastQueryParam(testURL)
			return []interface{}{last}, err
		},

		func(results []interface{}, t *testing.T) {
			last := results[0].(bool)
			if values.expectedLast != last {
				t.Fatalf("The value of 'last' should be '%v', got '%v'", values.expectedLast, last)
			}
		},
	)
}

func TestReadLastQueryParam(t *testing.T) {
	testValues := []readLastQueryParamTestValues{
		// intended usage
		{test: testValues{name: "NoParam", url: "http://localhost:9137/test", errorExpected: false}, expectedLast: false},
		{test: testValues{name: "ParamWithoutValue", url: "http://localhost:9137/test?last", errorExpected: false}, expectedLast: true},
		{test: testValues{name: "ParamTrue", url: "http://localhost:9137/test?last=1", errorExpected: false}, expectedLast: true},
		{test: testValues{name: "ParamFalse", url: "http://localhost:9137/test?last=False", errorExpected: false}, expectedLast: false},

		// expected errors
		{test: testValues{name: "ParamParseError", url: "http://localhost:9137/test?last=Fasle", errorExpected: true}},
		{test: testValues{name: "UnexpectedParam", url: "http://localhost:9137/test?last=True&other=something", errorExpected: true}},
	}

	for _, values := range testValues {
		t.Run(values.test.name, values.getTestFunc())
	}
}

type readSizeAndFromQueryParamsTestValues struct {
	test         testValues
	defaultSize  int
	expectedSize int
	expectedFrom string
}

func (values *readSizeAndFromQueryParamsTestValues) getTestFunc() func(*testing.T) {
	return values.test.getTestFunc(
		func(testURL *url.URL) ([]interface{}, error) {
			size, from, err := readSizeAndFromQueryParams(testURL, values.defaultSize)
			return []interface{}{size, from}, err
		},

		func(results []interface{}, t *testing.T) {
			size := results[0].(int)
			if values.expectedSize != size {
				t.Fatalf("The value of 'size' should be '%v', got '%v'", values.expectedSize, size)
			}

			from := results[1].(string)
			if values.expectedFrom != from {
				t.Fatalf("The value of 'from' should be '%v', got '%v'", values.expectedSize, from)
			}
		},
	)
}

func TestReadSizeAndFromQueryParam(t *testing.T) {
	testValues := []readSizeAndFromQueryParamsTestValues{
		// intended usage
		{test: testValues{name: "NoParam", url: "http://localhost:9137/test", errorExpected: false},
			defaultSize: 256, expectedSize: 256, expectedFrom: ""},
		{test: testValues{name: "SizeParam", url: "http://localhost:9137/test?size=512", errorExpected: false},
			defaultSize: 256, expectedSize: 512, expectedFrom: ""},
		{test: testValues{name: "FromParam", url: "http://localhost:9137/test?from=someTestId", errorExpected: false},
			defaultSize: 128, expectedSize: 128, expectedFrom: "someTestId"},
		{test: testValues{name: "BothParam", url: "http://localhost:9137/test?from=otherTestId&size=64", errorExpected: false},
			defaultSize: 512, expectedSize: 64, expectedFrom: "otherTestId"},

		// expected errors
		{test: testValues{name: "ParamWithoutValue", url: "http://localhost:9137/test?from", errorExpected: true}, defaultSize: 128},
		{test: testValues{name: "MultipleParam", url: "http://localhost:9137/test?size=16&from=objectId&size=32", errorExpected: true}, defaultSize: 128},
		{test: testValues{name: "ParamParseError", url: "http://localhost:9137/test?size=1ö24", errorExpected: true}, defaultSize: 64},
		{test: testValues{name: "UnexpectedParam", url: "http://localhost:9137/test?from=someId&other=something", errorExpected: true}, defaultSize: 64},
	}

	for _, values := range testValues {
		t.Run(values.test.name, values.getTestFunc())
	}
}
