package tensors

import (
	"github.com/sharnoff/testdep"
	"reflect"
	"testing"
)

func TestAll(t *testing.T) {
	g := testdep.New()

	// interpreter_test.go functions
	g.Require(tEquals, tNewInterpreter)
	g.Require(tCheckPoint, tNewInterpreter)
	g.Require(tCheckIndex, tNewInterpreter)
	g.Require(tIndex, tCheckPoint)
	g.Require(tPoint, tCheckIndex)
	g.Require(tIncrement, tCheckPoint)
	g.Require(tDecrement, tCheckPoint)
	g.Require(tIncreaseBy, tIncrement, tDecrement, tIndex, tPoint)

	// mapapply_test.go
	g.Require(tMapApply, tIncreaseBy, tIncrement)

	g.NameAll([]struct {
		Fn   func(*testing.T)
		Name string
	}{
		{tNewInterpreter, "NewInterpreter"},
		{tEquals, "Equals"},
		{tCheckPoint, "CheckPoint"},
		{tCheckIndex, "CheckIndex"},
		{tIndex, "Index"},
		{tPoint, "Point"},
		{tIncrement, "Increment"},
		{tDecrement, "Decrement"},
		{tIncreaseBy, "IncreaseBy"},
		{tMapApply, "MapApply"},
	})

	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}

	g.Test(t)
}

// returns true if there are no errors
func handleErrors(t *testing.T, name string, expected, got error, format string, a ...interface{}) bool {
	if expected == nil && got == nil {
		return true
	}

	args := make([]interface{}, len(a), len(a)+2) // plus two for each of the errors
	copy(args, a)
	args = append(args, expected, got)

	if format != "" {
		format += " "
	}
	format += "Expected %q, Got %q."

	if (expected == nil) != (got == nil) {
		t.Errorf(name+": Whether or not error was returned did not match. "+format, args...)
	} else if !Is(got, expected) {
		t.Errorf(name+": Unexpected error type. "+format, args...)
	}

	return false
}

// returns true if the returns are equal
func handleReturn(t *testing.T, name string, expected, got interface{}, format string, a ...interface{}) bool {
	if reflect.DeepEqual(got, expected) {
		return true
	}

	args := make([]interface{}, len(a), len(a)+2) // plus two for each of the errors
	copy(args, a)
	args = append(args, expected, got)

	if format != "" {
		format += " "
	}
	format += "Expected %v, Got %v."

	t.Errorf(name+": Bad return. "+format, args...)
	return false
}
