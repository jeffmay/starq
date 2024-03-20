package tsq

import (
	"fmt"
	"io"
)

// TODO: Add a builder interface, like so...
// type AnyBuilder[V AnyValue] interface {
//   String(string) V
//   Number(any) V
//   Object(map[string]V) V
// }

// TODO: Add MustExist, MustNotExist, MustGetNull (?)
type AnyValue interface {

	// Get returns the value at the given path or reports an error to the test runner.
	MustGet(path ...string) AnyValue

	// MustGetArray returns the value at the given path as an array or reports an error to the test runner.
	MustGetArray(path ...string) []AnyValue

	// MustGetObject returns the value at the given path as an object or reports an error to the test runner.
	MustGetObject(path ...string) map[string]AnyValue

	// MustGetString returns the value at the given path as a string or reports an error to the test runner.
	MustGetString(path ...string) string

	// MustGetInt64 returns the value at the given path as an int64 or reports an error to the test runner.
	MustGetInt64(path ...string) int64

	// MustGetFloat64 returns the value at the given path as a float64 or reports an error to the test runner.
	MustGetFloat64(path ...string) float64

	// MustGetBool returns the value at the given path as a bool or reports an error to the test runner.
	MustGetBool(path ...string) bool

	// Exists returns true if the value at the given path exists (may be null).
	Exists(path ...string) bool

	// IsNull returns true if the value at the given path exists and is null.
	IsNull(path ...string) bool

	// IsTop returns true if this value is the root of the document.
	IsTop() bool

	// IsEqual returns true if the given value is equal to this value.
	IsEqual(AnyValue) bool

	// PathFromTop returns the path from the root of the document to this value.
	PathFromTop() []string

	// Pretty returns a 2-space indented rendering of the document.
	//
	// Calling this multiple time MUST not require re-rendering the output.
	Pretty() string

	// Top returns the root of the document (which may be the same as the receiver).
	Top() AnyValue
}

type DocumentFormat string

const (
	JSON DocumentFormat = "json"
)

func ParseString(format DocumentFormat, data string) AnyValue {
	switch format {
	case JSON:
		return ParseJSON(data)
	default:
		panic(fmt.Errorf("unsupported format: %s", format))
	}
}

func Parse(format DocumentFormat, r io.Reader) (AnyValue, error) {
	return nil, nil
}
