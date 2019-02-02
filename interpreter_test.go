package tensors

import (
	"testing"
)

func tNewInterpreter(t *testing.T) {
	table := []struct {
		dims []int
		in   Interpreter
		err  error
	}{
		{[]int{1, 2, 3}, Interpreter{[]int{1, 2, 3}, []int{1, 2, 6}}, nil},
		{[]int{1, 2, 3, 4}, Interpreter{[]int{1, 2, 3, 4}, []int{1, 2, 6, 24}}, nil},
		{[]int{1, 2, 3, 4, 5}, Interpreter{[]int{1, 2, 3, 4, 5}, []int{1, 2, 6, 24, 120}}, nil},
		{[]int{2, 1, 3}, Interpreter{[]int{2, 1, 3}, []int{2, 2, 6}}, nil},
		{[]int{2, 3, 1}, Interpreter{[]int{2, 3, 1}, []int{2, 6, 6}}, nil},

		{nil, Interpreter{}, ErrZeroDims},
		{[]int{}, Interpreter{}, ErrZeroDims},
		{[]int{0, 1, 2}, Interpreter{}, DimsValueError{}},
		{[]int{1, 0, 2}, Interpreter{}, DimsValueError{}},
		{[]int{1, 2, 0}, Interpreter{}, DimsValueError{}},
	}

	for _, tab := range table {
		// check NewInterpreterSafe
		in, err := NewInterpreterSafe(tab.dims)

		_ = handleErrors(t, "NewInterpreter", tab.err, err, "Dims: %v.", tab.dims) &&
			handleReturn(t, "NewInterpreter", tab.in, in, "Dims: %v.", tab.dims)
	}
}

// requires NewInterpreter
func tEquals(t *testing.T) {
	table := []struct{
		a, b []int
		equals bool
	}{
		{[]int{1}, []int{1}, true},
		{[]int{5}, []int{5}, true},
		{[]int{1, 1, 1}, []int{1, 1, 1}, true},
		{[]int{6, 3, 8}, []int{6, 3, 8}, true},

		{[]int{1}, []int{1, 1}, false},
		{[]int{1, 1}, []int{1, 1, 1}, false},
		{[]int{1, 1, 1}, []int{1, 1, 2}, false},
		{[]int{6, 3, 8}, []int{8, 3, 6}, false},
	}

	for _, tab := range table {
		a, b := NewInterpreter(tab.a), NewInterpreter(tab.b)
		if Equals(a, b) != tab.equals {
			if tab.equals {
				t.Errorf("Equals: Expected %v = %v, got not equal.", a, b)
			} else {
				t.Errorf("Equals: Expected %v != %v, got equal.", a, b)
			}
		}
	}
}

// requires NewInterpreter
func tCheckPoint(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})

	table := []struct {
		point []int
		err   error
	}{
		{[]int{0, 0, 0}, nil},
		{[]int{1, 1, 3}, nil},
		{[]int{1, 2, 3}, nil},

		{nil, ErrZeroPoint},
		{[]int{}, ErrZeroPoint},

		{[]int{0, 0}, LengthMismatchError{}},
		{[]int{0, 0, 0, 0}, LengthMismatchError{}},

		{[]int{2, 0, 0}, PointOutOfBoundsError{}},
		{[]int{0, 0, 4}, PointOutOfBoundsError{}},
		{[]int{-1, 0, 0}, PointOutOfBoundsError{}},
		{[]int{0, 0, -1}, PointOutOfBoundsError{}},
	}

	for _, tab := range table {
		err := in.CheckPoint(tab.point)

		handleErrors(t, "CheckPoint", tab.err, err, "Interpreter: %v, Point: %v.", in, tab.point)
	}
}

// requires NewInterpreter
func tCheckIndex(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		index int
		err   error
	}{
		{0, nil},
		{23, nil},

		{-1, ErrIndexZero},
		{24, ErrIndexSize},
		{25, ErrIndexSize},
	}

	for _, tab := range table {
		err := in.CheckIndex(tab.index)

		handleErrors(t, "CheckIndex", tab.err, err, "Interpreter: %v, Index: %v.", in, tab.index)
	}
}

// requires CheckPoint
func tIndex(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		point []int
		index int
		err   error
	}{
		{nil, 0, ErrZeroPoint},
		{[]int{}, 0, ErrZeroPoint},
		{[]int{0, 0, 0, 0}, 0, LengthMismatchError{}},
		{[]int{0, 0}, 0, LengthMismatchError{}},
		{[]int{-1, 0, 0}, 0, PointOutOfBoundsError{}},
		{[]int{2, 2, 3}, 0, PointOutOfBoundsError{}},
		{[]int{1, 2, 4}, 0, PointOutOfBoundsError{}},

		{[]int{0, 0, 0}, 0, nil},
		{[]int{0, 1, 2}, 14, nil},
		{[]int{1, 2, 3}, 23, nil},
	}

	for _, tab := range table {
		index, err := in.IndexSafe(tab.point)

		// if there are no errors:
		_ = handleErrors(t, "Index", tab.err, err, "Interpreter: %v, Point: %v.", in, tab.point) &&
			handleReturn(t, "Index", tab.index, index, "Interpreter: %v, Point: %v.", in, tab.point)
	}
}

// requires CheckIndex
func tPoint(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		index int
		point []int
		err   error
	}{
		{0, []int{0, 0, 0}, nil},
		{23, []int{1, 2, 3}, nil},

		{-1, nil, ErrIndexZero},
		{24, nil, ErrIndexSize},
		{25, nil, ErrIndexSize},
	}

	for _, tab := range table {
		point, err := in.PointSafe(tab.index)

		// if there are no errors:
		_ = handleErrors(t, "Point", tab.err, err, "Interpreter: %v, Index: %v.", in, tab.index) &&
			handleReturn(t, "Point", tab.point, point, "Interpreter: %v, Index: %v.", in, tab.index)
	}
}

// requires CheckPoint
func tIncrement(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		point []int
		res   []int
		over  bool
		err   error
	}{
		{[]int{0, 0, 0}, []int{1, 0, 0}, true, nil},
		{[]int{1, 0, 0}, []int{0, 1, 0}, true, nil},
		{[]int{1, 2, 3}, nil, false, nil},

		{nil, nil, false, ErrZeroPoint},
		{[]int{}, nil, false, ErrZeroPoint},

		{[]int{0, 0}, nil, false, LengthMismatchError{}},
		{[]int{0, 0, 0, 0}, nil, false, LengthMismatchError{}},

		{[]int{2, 0, 0}, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, 4}, nil, false, PointOutOfBoundsError{}},
		{[]int{-1, 0, 0}, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, -1}, nil, false, PointOutOfBoundsError{}},
	}

	for _, tab := range table {
		point := make([]int, len(tab.point))
		copy(point, tab.point)

		over, err := in.IncrementSafe(point)

		format := "Interpreter: %v, Point: %v."
		a := []interface{}{in, tab.point}

		if handleErrors(t, "Increment", tab.err, err, format, a...) {
			if !tab.over || !over {
				handleReturn(t, "Increment", tab.over, over, format, a...)
			} else {
				handleReturn(t, "Increment", tab.res, point, format, a...)
			}
		}
	}
}

// requires CheckPoint
func tDecrement(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		point []int
		res   []int
		over  bool
		err   error
	}{
		{[]int{0, 0, 0}, nil, false, nil},
		{[]int{1, 0, 0}, []int{0, 0, 0}, true, nil},
		{[]int{0, 1, 0}, []int{1, 0, 0}, true, nil},

		{nil, nil, false, ErrZeroPoint},
		{[]int{}, nil, false, ErrZeroPoint},

		{[]int{0, 0}, nil, false, LengthMismatchError{}},
		{[]int{0, 0, 0, 0}, nil, false, LengthMismatchError{}},

		{[]int{2, 0, 0}, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, 4}, nil, false, PointOutOfBoundsError{}},
		{[]int{-1, 0, 0}, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, -1}, nil, false, PointOutOfBoundsError{}},
	}

	for _, tab := range table {
		point := make([]int, len(tab.point))
		copy(point, tab.point)

		over, err := in.DecrementSafe(point)

		format := "Interpreter: %v, Point: %v."
		a := []interface{}{in, tab.point}

		if handleErrors(t, "Decrement", tab.err, err, format, a...) {
			if !tab.over || !over {
				handleReturn(t, "Decrement", tab.over, over, format, a...)
			} else {
				handleReturn(t, "Decrement", tab.res, point, format, a...)
			}
		}
	}
}

// requires Increment, Decrement, Index, Point
func tIncreaseBy(t *testing.T) {
	in := NewInterpreter([]int{2, 3, 4})
	// size: 2*3*4 = 24

	table := []struct {
		point []int
		inc   int

		res  []int
		over bool
		err  error
	}{
		{[]int{0, 0, 0}, 0, []int{0, 0, 0}, true, nil},
		{[]int{1, 2, 3}, 0, []int{1, 2, 3}, true, nil},

		{[]int{0, 0, 0}, 3, []int{1, 1, 0}, true, nil},
		{[]int{0, 0, 0}, 23, []int{1, 2, 3}, true, nil},
		{[]int{1, 2, 3}, -3, []int{0, 1, 3}, true, nil},
		{[]int{1, 2, 3}, -23, []int{0, 0, 0}, true, nil},

		{[]int{0, 0, 0}, -1, nil, false, nil},
		{[]int{0, 0, 1}, -10, nil, false, nil},
		{[]int{1, 2, 3}, 1, nil, false, nil},
		{[]int{0, 0, 3}, 10, nil, false, nil},

		{nil, 1, nil, false, ErrZeroPoint},
		{[]int{}, 1, nil, false, ErrZeroPoint},

		{[]int{0, 0}, 1, nil, false, LengthMismatchError{}},
		{[]int{0, 0, 0, 0}, 1, nil, false, LengthMismatchError{}},

		{[]int{2, 0, 0}, 1, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, 4}, 1, nil, false, PointOutOfBoundsError{}},
		{[]int{-1, 0, 0}, 1, nil, false, PointOutOfBoundsError{}},
		{[]int{0, 0, -1}, 1, nil, false, PointOutOfBoundsError{}},

		{[]int{0, 0, 0}, 24, nil, false, ErrChangeTooBig},
		{[]int{0, 0, 0}, 25, nil, false, ErrChangeTooBig},
		{[]int{1, 2, 3}, -24, nil, false, ErrChangeTooBig},
		{[]int{1, 2, 3}, -25, nil, false, ErrChangeTooBig},
	}

	for _, tab := range table {
		point := make([]int, len(tab.point))
		copy(point, tab.point)

		over, err := in.IncreaseBySafe(point, tab.inc)

		format := "Interpreter: %v, Point: %v, Increase: %v."
		a := []interface{}{in, tab.point, tab.inc}

		if handleErrors(t, "IncreaseBySafe", tab.err, err, format, a...) {
			if !tab.over || !over {
				handleReturn(t, "IncreaseBySafe", tab.over, over, format, a...)
			} else {
				handleReturn(t, "IncreaseBySafe", tab.res, point, format, a...)
			}
		}
	}
}
