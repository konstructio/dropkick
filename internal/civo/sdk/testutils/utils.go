package testutils

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

// AssertEqual is a helper function to compare two values.
func AssertEqual[T comparable](t *testing.T, v1, v2 T) {
	t.Helper()

	if v1 != v2 {
		t.Fatalf("expected \"%v\", got \"%v\"", v1, v2)
	}
}

// AssertEqualf is a helper function to compare two values with a custom format.
func AssertEqualf[T comparable](t *testing.T, v1, v2 T, format string, args ...interface{}) {
	t.Helper()

	if v1 != v2 {
		t.Fatalf(format, args...)
	}
}

// AssertNoError is a helper function to check if an error is nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// AssertNoErrorf is a helper function to check if an error is nil with a custom format.
func AssertNoErrorf(t *testing.T, err error, format string, args ...interface{}) {
	t.Helper()

	if err != nil {
		t.Fatalf(format, args...)
	}
}

// AssertError is a helper function to check if an error is not nil.
func AssertErrorf(t *testing.T, err error, format string, args ...interface{}) {
	t.Helper()

	if err == nil {
		t.Fatalf(format, args...)
	}
}

// AssertErrorEqual is a helper function to compare two errors.
func AssertErrorEqual(t *testing.T, expected error, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		t.Fatalf("expected error to be \"%#v\", got \"%#v\"", expected, got)
	}
}

// MockCivo is a mock implementation of the Civoer interface.
type MockCivo struct {
	FnDo          func(ctx context.Context, location, method string, output interface{}, params map[string]string) error
	FnGetEndpoint func() string
	FnGetRegion   func() string
	FnGetClient   func() *http.Client
}

// Do is a mock implementation of the Civoer interface.
func (m *MockCivo) Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
	if m.FnDo == nil {
		return errors.New("method \"do\" not implemented")
	}

	return m.FnDo(ctx, location, method, output, params)
}

// GetClient is a mock implementation of the Civoer interface.
func (m *MockCivo) GetClient() *http.Client {
	if m.FnGetClient == nil {
		return nil
	}

	return m.FnGetClient()
}

// GetRegion is a mock implementation of the Civoer interface.
func (m *MockCivo) GetRegion() string {
	if m.FnGetRegion == nil {
		return ""
	}

	return m.FnGetRegion()
}

// GetEndpoint is a mock implementation of the Civoer interface.
func (m *MockCivo) GetEndpoint() string {
	if m.FnGetEndpoint == nil {
		return ""
	}

	return m.FnGetEndpoint()
}

// ParseParamNumber is a helper function to parse a number from a map of parameters.
// It will return the value of the key as an integer. If the key is not found, it will
// return 1. If the key is found but the value is not a number, it will fail the test.
// If shouldFail is set to true, the function will fail the test if the key is not found.
func ParseParamNumber(t *testing.T, location string, params map[string]string, key string, defval int, shouldFail bool) int {
	t.Helper()

	page, ok := params[key]
	if !ok {
		if shouldFail {
			t.Fatalf("expected page to be set for endpoint: %q", location)
		}

		// Default to page 1 if not found.
		return defval
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		t.Fatalf("failed to parse page number for endpoint %q: %v", location, err)
	}

	return pageNumber
}

// InjectIntoSlice is a helper function to inject a slice into another slice.
// The destination must be a pointer to a slice, and the source must be a slice.
func InjectIntoSlice(t *testing.T, destination, source interface{}) {
	t.Helper()

	// source must be a slice of some kind
	sliceValue := reflect.ValueOf(destination)

	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		t.Fatalf("expected a pointer to the object slice, got %v", sliceValue.Kind())
	}

	// destination needs to be a slice (doesn't matter what kind)
	sourceValue := reflect.ValueOf(source)

	if sourceValue.Kind() != reflect.Slice {
		t.Fatalf("expected a slice, got %v", sourceValue.Kind())
	}

	// append the source slice to the destination slice
	sliceValue.Elem().Set(reflect.AppendSlice(sliceValue.Elem(), sourceValue))
}
