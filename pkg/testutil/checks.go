package testutil

import (
	"reflect"
)

func (t *T) CheckTypeEquality(expected, actual interface{}) {
	t.Helper()

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	if expectedType != actualType {
		t.Errorf("Types do not match. Expected %s, Actual %s", expectedType, actualType)
		return
	}
}

func (t *T) CheckNoError(err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func (t *T) RequireNoError(err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error (failing test now): %s", err)
		t.FailNow()
	}
}

func (t *T) RequireNonNilResult(x interface{}, err error) interface{} {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error (failing test now): %s", err)
		t.FailNow()
	}
	if x == nil {
		t.Errorf("unexpected nil value (failing test now)")
		t.FailNow()
	}
	return x
}
