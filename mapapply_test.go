package tensors

import (
	"sync/atomic"
	"testing"
)

// requires Increment, IncreaseBy... etc
func tMapApply(t *testing.T) {
	in := NewInterpreter([]int{10, 15, 5})

	completed := make([]int64, in.Size())

	fn := func(point []int, index int) error {
		if in.Index(point) != index {
			t.Errorf("MapApply: fn given unequal point-index pair. Point: %v, Index: %v. in.Index(Point) = %v.",
				point, index, in.Index(point))
		}

		atomic.AddInt64(&(completed[index]), 1)
		return nil
	}

	threadOps := ThreadingOptions{10, 5}

	if err := in.MapApplySafe(fn, &threadOps); err != nil {
		t.Errorf("MapApply: Error returned when none expected. Got: %q.", err)
	}

	for i, c := range completed {
		if c != 1 {
			t.Errorf("MapApply: Index %d was not run once. Was run %d times.", i, c)
		}
	}
}
