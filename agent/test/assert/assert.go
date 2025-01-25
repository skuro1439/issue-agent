package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equal[V comparable](t *testing.T, got, expected V) {
	t.Helper()

	gotIsPtr := reflect.TypeOf(got).Kind() == reflect.Pointer
	expectedIsPtr := reflect.TypeOf(expected).Kind() == reflect.Pointer

	// compare pointers
	if gotIsPtr && expectedIsPtr {
		fmt.Println("both are pointers")
		t.Log("both are pointers")
		fmt.Printf("got: %v\n", reflect.ValueOf(got).Elem())
		if !reflect.ValueOf(got).Elem().Equal(reflect.ValueOf(expected).Elem()) {
			t.Errorf(`assert.Equal(
got: %v
expected: %v
)`, got, expected)
		}
		return
	}

	if expected != got {
		t.Errorf(`assert.Equal(
got: %v
expected: %v
)`, got, expected)
	}
}

func Nil(t *testing.T, value any) {
	t.Helper()

	if value != nil {
		t.Errorf(`execpted nil, got %v`, value)
	}
}

func EqualStringSlices(t *testing.T, got, expected []string) {
	t.Helper()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf(`assert.EqualStringSlices(
got: %v
expected: %v
)`, got, expected)
	}
}

func HasError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf(`expected error, got nil`)
	}
}

func NoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf(`expected no error, got %v`, err)
	}
}
