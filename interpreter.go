package tensors

// Interpreter is a general wrapper for indexing in n-dimensional slices. It is a central element
// of Tensors.
type Interpreter struct {
	// Dims makes up the sizes of each dimension. Eg: [width, height, depth, etc..]
	// Values are interpreted such that incrementing indices in Values increases (in order) the
	// indices of Dims[0], Dims[1], ... Dims[N]. In other words, changes in the indices of
	// higher-order dimensions result in greater changes in the base index.
	//
	// Dims should not be altered - it is set at construction. It is made public to be visible to
	// marshallers. However, it can (and should) be accessed if need be.
	Dims []int

	// Sizes stores the amount of values encapsulated by a 'set' of this dimension. To clarify:
	// Sizes[0] = Dims[0]; Sizes[1] = Dims[0]*Dims[1]; Sizes[N] = len(Values).
	//
	// Sizes should not be altered - it is set at construction. It is made public to be visible to
	// marshallers.
	Sizes []int
}

// NewInterpreter is fairly self-explanatory; it returns a new Interpreter, based on the provided
// dimensions.
//
// If the error conditions outlined in SafeInterpreter are not met, NewInterpreter will panic.

// Note: If any of the dimensions are <= 0 or len(dims) == 0, NewInterpreter will panic. This can
// be avoided with SafeInterpreter, which returns error.
//
// An additional note: NewInterpreter does NOT make a copy of dims -- if the array that dims
// references is modified, so too will dims.
func NewInterpreter(dims []int) Interpreter {
	in, err := NewInterpreterSafe(dims)
	if err != nil {
		panic(err)
	}

	return in
}

// SafeInterpreter returns a new Interpreter or any error. SafeInterpreter will return error if
// any one of the provided dimensions are <= 1, or if len(dims) == 0. These errors will either be
// ErrZeroDims or a DimsValueError.
//
// Note: SafeInterpreter does NOT make a copy of dims -- if the array that dims references is
// modified, dims will be also.
func NewInterpreterSafe(dims []int) (Interpreter, error) {
	if len(dims) == 0 {
		return Interpreter{}, ErrZeroDims
	} else {
		for i, d := range dims {
			if d <= 0 {
				return Interpreter{}, DimsValueError{dims, i}
			}
		}
	}

	sizes := make([]int, len(dims))
	sizes[0] = dims[0]

	for i := 1; i < len(sizes); i++ {
		sizes[i] = sizes[i-1] * dims[i]
	}

	return Interpreter{dims, sizes}, nil
}

// CheckPoint is mostly for internal use. It checks that the point is within the space defined by
// the Interpreter and returns error if it is not. CheckPoint is made public to allow for a common
// place to define constraints on points and the expected behavior if those constraints are not
// kept.
//
// CheckPoint has three (and a half) error conditions:
//		(0) If the length of point is zero.
//		(1)	If the number of dimensions given by point is not equal to the number of dimensions of
//			the Interpreter.
//		(2)	If any of the values of point are below zero
//		(3) If any of the values of point are greater than their respective dimension sizes,
//			according to the Interpreter. In other words, if point[i] ≥ in.Dims[i] for any i.
// (0) will return an ErrZeroPoint, (1) will return a LengthMismatchError, (2) and (3) will return
// PointOutOfBoundsError.
func (in Interpreter) CheckPoint(point []int) error {
	if len(point) == 0 {
		return ErrZeroPoint
	} else if len(point) != len(in.Sizes) {
		return LengthMismatchError{"point", len(point), len(in.Sizes)}
	}

	for i, v := range point {
		if v < 0 || v >= in.Dims[i] {
			return PointOutOfBoundsError{point, in.Dims, i}
		}
	}

	return nil
}

// Index returns the index in the base array that the given point corresponds to. Index will panic
// if any of the criteria documented by Interpreter.CheckPoint() are not met.
func (in Interpreter) Index(point []int) int {
	index, err := in.IndexSafe(point)
	if err != nil {
		panic(err)
	}

	return index
}

// IndexSafe is a 'safe' version of Index. It will return any errors that may occur. Possible
// errors can be found in the documentation for Interpreter.CheckPoint().
func (in Interpreter) IndexSafe(point []int) (int, error) {
	if err := in.CheckPoint(point); err != nil {
		return 0, err
	}

	return in.IndexFast(point), nil
}

// IndexFast is the 'Fast' variant of Index. It does not check for any error conditions. If the
// error conditions have already been checked, IndexFast will not panic.
//
// IndexFast exists for those occasions when you know that the contidions have already been met.
// Please note: IndexFast is only marginally faster than Index; it is simply avoiding a call to
// Interpreter.CheckPoint().
func (in Interpreter) IndexFast(point []int) int {
	index := point[0]
	for i := 1; i < len(in.Sizes); i++ {
		index += point[i] * in.Sizes[i-1]
	}

	return index
}

// CheckIndex is mostly for internal use. It checks that the given index is within the bounds of
// the Interpreter's size. CheckIndex is made public to allow for a common place to define
// constraints on indices and give the expected behavior if those constraints are not kept.
//
// CheckIndex has two error conditions:
//		(0) If index is less than 0
//		(1) If index ≥ the Size of the Interpreter
// (0) will return ErrIndexZero and (1) will return ErrIndexSize.
func (in Interpreter) CheckIndex(index int) error {
	if index < 0 {
		return ErrIndexZero
	} else if index >= in.Size() {
		return ErrIndexSize
	}

	return nil
}

// Point returns the multi-dimensional point corresponding to the given index in the base array.
// If any of the errors that would be encountered by Interpreter.CheckIndex() are encountered here,
// Point will panic with that error.
func (in Interpreter) Point(index int) []int {
	p, err := in.PointSafe(index)
	if err != nil {
		panic(err)
	}

	return p
}

// PointSafe is the base operation that is run by Point. Where Point panics, PointSafe returns
// error. Possible errors can be found in the documentation ofr Interpreter.CheckIndex().
func (in Interpreter) PointSafe(index int) ([]int, error) {
	if err := in.CheckIndex(index); err != nil {
		return nil, err
	}

	p := make([]int, len(in.Dims))
	for i := len(p) - 1; i >= 1; i-- {
		p[i] = index / in.Sizes[i-1]
		index %= in.Sizes[i-1]
	}

	p[0] = index
	return p, nil
}

// Size returns the required (and expected) length of the base array for the Interpreter.
func (in Interpreter) Size() int {
	// the size is equal to the size of the largest dimension.
	return in.Sizes[len(in.Sizes)-1]
}

// Increment increases the index corresponding to the point by 1. If the point is already at the
// maximum index (so incrementing would overflow), Increment returns false. Otherwise, it returns
// true. Additionally, Increment WILL NOT preserve the original value of the point if it is at its
// maximum value -- it will overflow to a zero'd point.
//
// If point does not meet the necessary conditions mappped out in the documentation for CheckIndex,
// Increment will panic with those errors.
func (in Interpreter) Increment(point []int) bool {
	b, err := in.IncrementSafe(point)
	if err != nil {
		panic(err)
	}

	return b
}

// IncrementSafe undergoes the same operation as Increment, but returns error instead of panicking.
// IncrementSafe will return error if any of the conditions outlined in the documentation for
// CheckIndex aren't met.
//
// If IncrementSafe returns error, the original point has not changed.
func (in Interpreter) IncrementSafe(point []int) (bool, error) {
	if err := in.CheckPoint(point); err != nil {
		return false, err
	}

	return in.IncrementFast(point), nil
}

// IncrementFast performs the same operation as Increment, but does not check for possible error
// conditions. As such, possible errors can be encountered from IncrementFast that will not be
// visible until later.
//
// IncrementFast makes debugging your code harder.
func (in Interpreter) IncrementFast(point []int) bool {
	for i := 0; i < len(in.Dims); i++ {
		point[i]++
		if point[i] < in.Dims[i] {
			break
		}

		point[i] = 0

		if i == len(in.Dims)-1 {
			return false
		}
	}

	return true
}

// Decrement is the decreasing analog to Increment. If point is at its mimimum value, Decrement
// returns false and shifts point to its maximum value as an overflow. Increment requires the same
// conditions of point as Increment does.
//
// Assume Decrement (and its derivatives) behaves in the same way as Increment (and its
// derivatives). Decrement panics where DecrementSafe returns error.
func (in Interpreter) Decrement(point []int) bool {
	b, err := in.DecrementSafe(point)
	if err != nil {
		panic(err)
	}

	return b
}

// DecrementSafe is the decreasing analog to IncrementSafe. DecrementSafe returns error if
// Interpreter.CheckPoint() would return eror.
func (in Interpreter) DecrementSafe(point []int) (bool, error) {
	if err := in.CheckPoint(point); err != nil {
		return false, err
	}

	return in.DecrementFast(point), nil
}

// DecrementFast is the decreasing analog to IncrementFast.
func (in Interpreter) DecrementFast(point []int) bool {
	for i := 0; i < len(in.Dims); i++ {
		point[i]--
		if point[i] >= 0 {
			break
		}

		point[i] = in.Dims[i] - 1

		if i == len(in.Dims)-1 {
			return false
		}
	}

	return true
}

// IncreaseBy increases the index corresponding to the point by the given value. IncreaseBy will
// panic if any of the following conditions are not met: (0) point must be 'valid' -- see
// Interpreter.Index() and (1) change must be >= 0 and < Interpreter.Size(). Additionally, if the
// magnitude of the change is larger than the size of the Interpreter, Increase will panic with
// ErrChangeTooBig.
//
// If the increase to point would overflow it, IncreaseBy returns false (in a similar fashion to
// Increment). In those cases, the changes to point are not defined -- point should not be used
// after IncreaseBy returns 'false'.
//
// IncreaseBy can be given negative values. If given a change of 0, IncreaseBy will always return
// (true, nil). Additionally, if change equals ±1, Increment and Decrement will be used instead.*
//
// * It will actually be IncrementSafe and DecrementSafe, because IncreaseBy uses IncreaseBySafe
// internally.
func (in Interpreter) IncreaseBy(point []int, change int) bool {
	b, err := in.IncreaseBySafe(point, change)
	if err != nil {
		panic(err)
	}

	return b
}

// IncreaseBySafe performs the same operation as IncreaseBy, but will return error instead of
// panicking.
func (in Interpreter) IncreaseBySafe(point []int, change int) (bool, error) {
	if change == 0 {
		return true, nil
	} else if change == 1 {
		return in.IncrementSafe(point)
	} else if change == -1 {
		return in.DecrementSafe(point)
	}

	index, err := in.IndexSafe(point)
	if err != nil {
		return false, err
	}

	if change >= in.Size() {
		return false, ErrChangeTooBig
	} else if change <= -in.Size() {
		return false, ErrChangeTooBig
	}

	newIndex := index + change

	if newIndex < 0 || newIndex >= in.Size() {
		return false, nil
	}

	newPoint := in.Point(newIndex)
	copy(point, newPoint)

	return true, nil
}
