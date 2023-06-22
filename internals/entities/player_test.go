package entities_test

import (
	"testing"

	"github.com/ygortgaleno/cloudwalk-test/internals/entities"
)

type updatePayerKillTest struct {
	initialValue, expected int
}

var upkTests = []updatePayerKillTest{
	{2, 3},
	{4, 5},
	{20, 21},
	{97, 98},
}

type updateWorldDeathTest struct {
	initialValue, givenValue, expected int
}

var uwdTests = []updateWorldDeathTest{
	{2, 3, 5},
	{4, 5, 9},
	{20, 21, 41},
	{9, 1, 10},
}

func TestUpdatePlayerKill(t *testing.T) {
	for _, test := range upkTests {
		p := entities.Player{Name: "player", Kills: uint32(test.initialValue), WorldDeaths: 0}
		p.UpdatePlayerKill()

		expected := test.expected
		received := p.Kills

		if received != uint32(expected) {
			t.Errorf("expected %q, received %q", expected, received)
		}
	}
}

func TestUpdateWorldDeath(t *testing.T) {
	for _, test := range uwdTests {
		p := entities.Player{Name: "player", Kills: 0, WorldDeaths: int64(test.initialValue)}
		p.UpdateWorldDeaths(int64(test.givenValue))

		expected := test.expected
		received := p.WorldDeaths

		if received != int64(expected) {
			t.Errorf("expected %q, received %q", expected, received)
		}
	}
}
