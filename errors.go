package tensors

import (
	"fmt"
	"reflect"
)

// DimsValueError serves to document errors from dimensions (not point values) being non-positive.
type DimsValueError struct {
	dims  []int
	index int
}

// LengthMismatchError serves to document errors having to do with lengths of provided points or
// dimensions. For example, if the point supplied to Interpreter.Index() was too short.
type LengthMismatchError struct {
	variant string

	is       int
	shouldBe int
}

// PointOutOfBoundsError serves to document errors having to do with indices of points being out of
// the bounds of the dimensions, either less than 0 or greater than that index of the dimensions.
type PointOutOfBoundsError struct {
	point []int
	dims  []int
	index int
}

func (err DimsValueError) Error() string {
	return fmt.Sprintf("dims[%d] ≤ 0. dims: %v", err.index, err.dims)
}

func (err LengthMismatchError) Error() string {
	return fmt.Sprintf(err.variant+" length mismatch (is: %d, should be: %d)", err.is, err.shouldBe)
}

func (err PointOutOfBoundsError) Error() string {
	return fmt.Sprintf("point[%d] = %d is out of bounds of dims[%d] = %d",
		err.index, err.point[err.index], err.index, err.dims[err.index])
}

// Is checks whether or not two errors from this package are the same type. This is more than just
// a simple type comparison; Is checks whether or not the errors are, fundamentally, the same
// error. For type tensors.Error, Is checks individual variables (eg. ErrZeroDims != ErrZeroPoint),
// and for other types (eg. DimsValueError and LengthMismatchError) Is performs a type comparison.
//
// Is uses reflect, so it should only be run when an error has actually occurred.
func Is(err, base error) bool {
	if e0, ok := base.(Error); ok {
		if e1, ok := err.(Error); ok {
			return e0 == e1
		}

		return false
	}

	return reflect.TypeOf(base) == reflect.TypeOf(err)
}

// Error is a placeholder for specific errors and their messages, all of which are stored as vars
type Error struct{ string }

func (e Error) Error() string { return e.string }

var (
	ErrZeroDims       = Error{"dims has len = 0"}
	ErrZeroPoint      = Error{"point has len = 0"}
	ErrIndexZero      = Error{"index is < 0"}
	ErrIndexSize      = Error{"index is greater than Interpreter size"}
	ErrChangeTooBig   = Error{"magnitude of change is greater than Interpreter Size"}
	ErrPointOutOfSync = Error{"increasing point failed while index was within bounds"}
	ErrNilFunction    = Error{"given MapApply function is nil"}
)
