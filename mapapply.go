package tensors

import "sync"

// ThreadingOptions serves only as an argument to MapApply and its derivatives. It serves to group
// optional arguments for multithreading with MapApply.
type ThreadingOptions struct {
	// OpsPerThread determines how many times each thread will call the function given to MapApply.
	OpsPerThread int

	// NumThreads determines the number of goroutines that are created to call the function.
	NumThreads int
}

// MapApply applies a given function to every value, giving the point and the index corresponding
// to the current value. MapApply iterates with increasing indices over all values of the
// Interpreter. ThreadingOptions is given to configure the specifics on the ratios for
// multithreading. If options is nil, MapApply will run as a single thread.
//
// Information on types of errors that can be recovered from here is documented with MapApplySafe.
//
// Additionally, if individual members of options are less than 1, they will be set to 1. This
// means that fields in options that are not explicitly set will default to 1. However, 1 is not
// an optimal value for multithreading, so it is not recommended.
//
// If multithreaded, fn will not recieve copies of 'point', so it SHOULD NOT be modified.
func (in Interpreter) MapApply(fn func([]int, int), options *ThreadingOptions) {
	// trim the puzzle pieces so they fit
	newFn := func(point []int, index int) error {
		fn(point, index)
		return nil
	}

	if err := in.MapApplySafe(newFn, options); err != nil {
		panic(err)
	}

	return
}

// MapApplySafe is effectively the same as MapApply, except it will return error instead of
// panicking. It also differs slightly from MapApply and MapApplyFast in that it expects fn to
// return error.
//
// There are several categories of errors that MapApplySafe could return, only some of which are
// 'expected' errors generated originally at MapApply:
//		(0) ErrNilFunction will be returned if fn is nil
//		(1) Any error returned by 'fn' will be passed along without context
//		(2) ErrPointOutOfSync and errors of type LengthMismatchError and PointOutOfBoundsError may
//			also be returned, due to internal problems.
func (in Interpreter) MapApplySafe(fn func([]int, int) error, options *ThreadingOptions) error {
	return in.generalMapApply(fn, options, false)
}

// MapApplyFast is functionally the same as MapApply, but it uses the 'Fast' variants of other
// functions instead, such as IncrementFast. This is available based on the principle that many
// small time reductions may greatly decrease computation time.
//
// MapApplyFast will make your code harder to debug. However, if it works with MapApply, it will
// work with MapApplyFast.
func (in Interpreter) MapApplyFast(fn func([]int, int), options *ThreadingOptions) {
	newFn := func(point []int, index int) error {
		fn(point, index)
		return nil
	}

	// we do check for errors here, because we don't want to ignore them if they do happen.
	// However, errors should realistically NEVER happen here, because none of the functions
	// called return a non-nil error
	if err := in.generalMapApply(newFn, options, true); err != nil {
		panic(err)
	}

	return
}

// generalMapApply serves as a way to reduce repetition within the MapApply functions
// useFast indicates whether or not to use the 'Fast' variants of functions. If true, it will not
// return error.
//
// generalMapApply returns two original errors: ErrPointOutOfSync, when increasing the value of a
// point would overflow sooner than expected, and ErrNilFunction. Other errors come from Increment
// and IncreaseBy, in addition to the user-supplied function: fn
func (in Interpreter) generalMapApply(fn func([]int, int) error, options *ThreadingOptions, useFast bool) error {
	if !useFast && fn == nil {
		return ErrNilFunction
	}

	var increment func([]int) (bool, error)
	var makePoint func(int) ([]int, error)

	// set functions depending upon what type we're using
	if useFast {
		increment = func(point []int) (bool, error) {
			return in.IncrementFast(point), nil
		}

		makePoint = func(index int) ([]int, error) {
			return in.Point(index), nil
		}
	} else {
		increment = in.IncrementSafe
		makePoint = in.PointSafe
	}

	// fill in threading options.
	if options == nil {
		// opsPerThread equal to in.Size to avoid the overhead of repeatedly checking back to get more.
		// this could also work with (1, 1), but setting OpsPerThread equal to in.Size() is faster.
		options = &ThreadingOptions{NumThreads: 1, OpsPerThread: in.Size()}
	} else {
		if options.OpsPerThread < 1 {
			options.OpsPerThread = 1
		}
		if options.NumThreads < 1 {
			options.NumThreads = 1
		}
	}

	// define a helper function that we'll use later
	dupe := func(p []int) []int {
		newP := make([]int, len(p))
		copy(newP, p)
		return newP
	}

	var index int
	point := make([]int, len(in.Dims))
	var err error

	var mux sync.Mutex
	var wg sync.WaitGroup

	// removed by a defer statement at the top of the anonymous function
	wg.Add(options.NumThreads)
	for thread := 0; thread < options.NumThreads; thread++ {
		go func() {
			defer wg.Done()

			// this is where errors from this goroutine are reported, just so that our syntax is a
			// little cleaner - we get to avoid more calls to mux.Lock() and mux.Unlock()
			var localErr error

			// these are the local variables that we'll use to iterate over a set of inputs to fn
			var localIndex int
			var localPoint []int

			// localEnd is the last index (exclusive) of our range.
			var localEnd int

			// This outer loop obtains more values to iterate over.
			for {
				// get more values, check for errors
				mux.Lock()
				if err != nil {
					mux.Unlock()
					return
				} else if localErr != nil {
					err = localErr
					mux.Unlock()
					return
				}

				// If there are no more args to use, exit
				if index >= in.Size() {
					mux.Unlock()
					return
				}

				localIndex = index
				localPoint = dupe(point)

				index += options.OpsPerThread
				localEnd = index

				// we have this in a separate condition because we don't care about how the point
				// if we already know that we're done. This way, we can avoid false errors from
				// indices that would be out of bounds
				if index < in.Size() {
					if options.OpsPerThread > 1 {
						if point, err = makePoint(index); err != nil {
							mux.Unlock()
							return
						}
					} else {
						var cont bool
						cont, err = increment(point)
						if err != nil {
							mux.Unlock()
							return
						} else if !cont {
							err = ErrPointOutOfSync
							mux.Unlock()
							return
						}
					}

				}

				mux.Unlock()

				// avoid overflow. This is done here, and not while we have the mux lock because it
				// would be a waste of computation to do there. We don't need to hold the lock to
				// do it.
				if localEnd > in.Size() {
					localEnd = in.Size()
				}

				// loop through the args we've fetched
				for localIndex < localEnd {
					if err := fn(localPoint, localIndex); err != nil {
						localErr = err
					}

					localIndex++

					// we don't want to increment the point. It's both wasted computation and will
					// result in a false error -- at this point we've finished the set of inputs.
					if localIndex == localEnd {
						break
					}

					// these lines mirror those above, where we have possession of the mux lock.
					cont, err := increment(localPoint)
					if err != nil {
						// loop back to the top, where we can record this error
						localErr = err
						break
					} else if !cont {
						localErr = ErrPointOutOfSync
						break
					}
				}
			}
		}()
	}

	wg.Wait()

	// will return nil if everything's fine
	return err
}
