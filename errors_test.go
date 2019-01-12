package tensors

import "testing"

// tests the equality of every type of error that is returned by tensors, to ensure that they are
// all recognized as unique
func TestTypedErrors(t *testing.T) {
	errs := []error{
		DimsValueError{},
		LengthMismatchError{},
		PointOutOfBoundsError{},

		ErrZeroDims,
		ErrZeroPoint,
		ErrIndexZero,
		ErrIndexSize,
		ErrChangeTooBig,
		ErrPointOutOfSync,
		ErrNilFunction,
	}

	for i := range errs {
		for j := range errs {
			if Is(errs[i], errs[j]) != (i == j) {
				t.Errorf("Identifier func Is wrongly identified type %T==%T as %v",
					errs[i], errs[j], !(i == j))
			}
		}
	}
}
