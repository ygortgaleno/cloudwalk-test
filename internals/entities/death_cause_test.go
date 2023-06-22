package entities_test

import (
	"testing"

	"github.com/ygortgaleno/cloudwalk-test/internals/entities"
)

type updateCounterTest struct {
	initialValue, expected int
}

var ucTests = []updateCounterTest{
	{2, 3},
	{4, 5},
	{20, 21},
	{97, 98},
}

func TestUpdateCounter(t *testing.T) {
	for _, test := range ucTests {

		dc := entities.DeathCause{Name: "TEST_CAUSE", Counter: uint(test.initialValue)}
		dc.UpdateCounter()

		expected := test.expected
		received := dc.Counter

		if received != uint(expected) {
			t.Errorf("expected %q, received %q", expected, received)
		}
	}
}
