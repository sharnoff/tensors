package tensors

// Tensor is a general type for facilitating the use of mathematical tensors. They consist of a
// base location for the storage of data, in addition to the
type Tensor struct {
	Interpreter

	// Values is the set of the values stored in the Tensor. The description for the storage of
	// these values can be found in the documentation for Interpreter.Dims
	Values []float64
}

// NewTensor returns a new Tensor, and will panic if any of the error conditions from
// NewInterpreterSafe are met.
func NewTensor(dims []int) Tensor {
	in := NewInterpreter(dims)
	return Tensor{in, make([]float64, in.Size())}
}

// NewTensorSafe undergoes the same process as NewTensor, but returns error instead of panicking.
func NewTensorSafe(dims []int) (Tensor, error) {
	// delegate checking dims to interpreter construction
	in, err := NewInterpreterSafe(dims)
	if err != nil {
		return Tensor{}, err
	}

	return Tensor{in, make([]float64, in.Size())}, nil
}

// PointValue returns the value of the tensor at the given point. PointValue requires the same
// conditions as Interpreter.Index (and thus, Interpreter.CheckPoint).
func (t Tensor) PointValue(point []int) float64 {
	return t.Values[t.Index(point)]
}

// PointValueSafe undergoes the same process as PointValue, but will return error instead of
// panicking.
func (t Tensor) PointValueSafe(point []int) (float64, error) {
	index, err := t.IndexSafe(point)
	if err != nil {
		return 0.0, err
	}

	return t.Values[index], nil
}
