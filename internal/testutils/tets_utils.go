package testutils

import "testing"

// func NewTestApplication(t *testing.T) *application {

// }

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}
