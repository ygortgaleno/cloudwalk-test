package services_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/ygortgaleno/cloudwalk-test/internals/dtos"
	"github.com/ygortgaleno/cloudwalk-test/internals/services"
)

type quakeLogParserTest struct {
	expeted map[string]dtos.GameDto
}

var qlpTest = quakeLogParserTest{
	expeted: map[string]dtos.GameDto{
		"game_1": {
			TotalKills: 0, Players: []string{}, Kills: map[string]int64{}, KillsByMeans: map[string]uint{},
		},
		"game_2": {
			TotalKills: 11,
			Players:    []string{"Isgalamido", "Mocinha"},
			Kills: map[string]int64{
				"Isgalamido": -5,
				"Mocinha":    0,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       1,
				"MOD_ROCKET_SPLASH": 3,
				"MOD_TRIGGER_HURT":  7,
			},
		},
		"game_3": {
			TotalKills: 4,
			Players:    []string{"Isgalamido", "Mocinha", "Zeh"},
			Kills: map[string]int64{
				"Isgalamido": 1,
				"Mocinha":    -1,
				"Zeh":        -2,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":      1,
				"MOD_ROCKET":       1,
				"MOD_TRIGGER_HURT": 2,
			},
		},
		"game_4": {
			TotalKills: 105,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 13,
				"Dono da Bola":   13,
				"Isgalamido":     19,
				"Zeh":            20,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       11,
				"MOD_MACHINEGUN":    4,
				"MOD_RAILGUN":       8,
				"MOD_ROCKET":        20,
				"MOD_ROCKET_SPLASH": 51,
				"MOD_SHOTGUN":       2,
				"MOD_TRIGGER_HURT":  9,
			},
		},
		"game_5": {
			TotalKills: 14,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 1,
				"Dono da Bola":   1,
				"Isgalamido":     2,
				"Zeh":            0,
			},
			KillsByMeans: map[string]uint{
				"MOD_RAILGUN":       1,
				"MOD_ROCKET":        4,
				"MOD_ROCKET_SPLASH": 4,
				"MOD_TRIGGER_HURT":  5,
			},
		},
		"game_6": {
			TotalKills: 29,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Oootsimo", "UnnamedPlayer", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 1,
				"Dono da Bola":   2,
				"Isgalamido":     3,
				"Oootsimo":       8,
				"UnnamedPlayer":  0,
				"Zeh":            7,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       1,
				"MOD_MACHINEGUN":    1,
				"MOD_RAILGUN":       2,
				"MOD_ROCKET":        5,
				"MOD_ROCKET_SPLASH": 13,
				"MOD_SHOTGUN":       4,
				"MOD_TRIGGER_HURT":  3,
			},
		},
		"game_7": {
			TotalKills: 130,
			Players:    []string{"Assasinu Credi", "Chessus", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 22,
				"Chessus":        0,
				"Dono da Bola":   12,
				"Isgalamido":     16,
				"Mal":            -3,
				"Oootsimo":       20,
				"Zeh":            9,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       7,
				"MOD_MACHINEGUN":    9,
				"MOD_RAILGUN":       9,
				"MOD_ROCKET":        29,
				"MOD_ROCKET_SPLASH": 49,
				"MOD_SHOTGUN":       7,
				"MOD_TRIGGER_HURT":  20,
			},
		},
		"game_8": {
			TotalKills: 89,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 10,
				"Dono da Bola":   3,
				"Isgalamido":     20,
				"Mal":            -2,
				"Oootsimo":       16,
				"Zeh":            12,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       6,
				"MOD_MACHINEGUN":    4,
				"MOD_RAILGUN":       12,
				"MOD_ROCKET":        18,
				"MOD_ROCKET_SPLASH": 39,
				"MOD_SHOTGUN":       1,
				"MOD_TRIGGER_HURT":  9,
			},
		},
		"game_9": {
			TotalKills: 67,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 10,
				"Dono da Bola":   11,
				"Isgalamido":     0,
				"Mal":            3,
				"Oootsimo":       9,
				"Zeh":            12,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       3,
				"MOD_MACHINEGUN":    3,
				"MOD_RAILGUN":       10,
				"MOD_ROCKET":        17,
				"MOD_ROCKET_SPLASH": 25,
				"MOD_SHOTGUN":       1,
				"MOD_TRIGGER_HURT":  8,
			},
		},
		"game_10": {
			TotalKills: 60,
			Players:    []string{"Assasinu Credi", "Chessus", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 3,
				"Chessus":        5,
				"Dono da Bola":   3,
				"Isgalamido":     6,
				"Mal":            1,
				"Oootsimo":       -1,
				"Zeh":            7,
			},
			KillsByMeans: map[string]uint{
				"MOD_BFG":           2,
				"MOD_BFG_SPLASH":    2,
				"MOD_CRUSH":         1,
				"MOD_MACHINEGUN":    1,
				"MOD_RAILGUN":       7,
				"MOD_ROCKET":        4,
				"MOD_ROCKET_SPLASH": 1,
				"MOD_TELEFRAG":      25,
				"MOD_TRIGGER_HURT":  17,
			},
		},
		"game_11": {
			TotalKills: 20,
			Players:    []string{"Assasinu Credi", "Chessus", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": -3,
				"Chessus":        0,
				"Dono da Bola":   -2,
				"Isgalamido":     5,
				"Mal":            0,
				"Oootsimo":       4,
				"Zeh":            0,
			},
			KillsByMeans: map[string]uint{
				"MOD_BFG_SPLASH":    3,
				"MOD_CRUSH":         1,
				"MOD_MACHINEGUN":    1,
				"MOD_RAILGUN":       4,
				"MOD_ROCKET_SPLASH": 4,
				"MOD_TRIGGER_HURT":  7,
			},
		},
		"game_12": {
			TotalKills: 160,
			Players:    []string{"Assasinu Credi", "Chessus", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 20,
				"Chessus":        13,
				"Dono da Bola":   3,
				"Isgalamido":     26,
				"Mal":            -6,
				"Oootsimo":       13,
				"Zeh":            13,
			},
			KillsByMeans: map[string]uint{
				"MOD_BFG":           8,
				"MOD_BFG_SPLASH":    8,
				"MOD_FALLING":       2,
				"MOD_MACHINEGUN":    7,
				"MOD_RAILGUN":       38,
				"MOD_ROCKET":        25,
				"MOD_ROCKET_SPLASH": 35,
				"MOD_TRIGGER_HURT":  37,
			},
		},
		"game_13": {
			TotalKills: 6,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 0,
				"Dono da Bola":   -1,
				"Isgalamido":     -1,
				"Oootsimo":       2,
				"Zeh":            2,
			},
			KillsByMeans: map[string]uint{
				"MOD_BFG":           1,
				"MOD_BFG_SPLASH":    1,
				"MOD_ROCKET":        1,
				"MOD_ROCKET_SPLASH": 1,
				"MOD_TRIGGER_HURT":  2,
			},
		},
		"game_14": {
			TotalKills: 122,
			Players:    []string{"Assasinu Credi", "Chessus", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 7,
				"Chessus":        7,
				"Dono da Bola":   2,
				"Isgalamido":     22,
				"Mal":            -2,
				"Oootsimo":       9,
				"Zeh":            5,
			},
			KillsByMeans: map[string]uint{
				"MOD_BFG":           5,
				"MOD_BFG_SPLASH":    10,
				"MOD_FALLING":       5,
				"MOD_MACHINEGUN":    4,
				"MOD_RAILGUN":       20,
				"MOD_ROCKET":        23,
				"MOD_ROCKET_SPLASH": 24,
				"MOD_TRIGGER_HURT":  31,
			},
		},
		"game_15": {
			TotalKills: 3,
			Players:    []string{"Zeh"},
			Kills: map[string]int64{
				"Zeh": -3,
			},
			KillsByMeans: map[string]uint{
				"MOD_TRIGGER_HURT": 3,
			},
		},
		"game_16": {
			TotalKills: 0, Players: []string{}, Kills: map[string]int64{}, KillsByMeans: map[string]uint{},
		},
		"game_17": {
			TotalKills: 13,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": -3,
				"Dono da Bola":   -2,
				"Isgalamido":     0,
				"Mal":            -1,
				"Oootsimo":       1,
				"Zeh":            0,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       3,
				"MOD_RAILGUN":       2,
				"MOD_ROCKET_SPLASH": 2,
				"MOD_TRIGGER_HURT":  6,
			},
		},
		"game_18": {
			TotalKills: 7,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 2,
				"Dono da Bola":   -1,
				"Isgalamido":     1,
				"Mal":            -1,
				"Oootsimo":       0,
				"Zeh":            2,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       1,
				"MOD_ROCKET":        1,
				"MOD_ROCKET_SPLASH": 4,
				"MOD_TRIGGER_HURT":  1,
			},
		},
		"game_19": {
			TotalKills: 95,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Isgalamido", "Mal", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 9,
				"Dono da Bola":   14,
				"Isgalamido":     14,
				"Mal":            2,
				"Oootsimo":       10,
				"Zeh":            20,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       1,
				"MOD_MACHINEGUN":    7,
				"MOD_RAILGUN":       10,
				"MOD_ROCKET":        27,
				"MOD_ROCKET_SPLASH": 32,
				"MOD_SHOTGUN":       6,
				"MOD_TRIGGER_HURT":  12,
			},
		},
		"game_20": {
			TotalKills: 3,
			Players:    []string{"Assasinu Credi", "Dono da Bola", "Oootsimo", "Zeh"},
			Kills: map[string]int64{
				"Assasinu Credi": 0,
				"Dono da Bola":   2,
				"Oootsimo":       1,
				"Zeh":            0,
			},
			KillsByMeans: map[string]uint{
				"MOD_ROCKET":        1,
				"MOD_ROCKET_SPLASH": 2,
			},
		},
		"game_21": {
			TotalKills: 131,
			Players:    []string{"Zeh", "Isgalamido", "Dono da Bola", "Assasinu Credi", "Mal", "Oootsimo"},
			Kills: map[string]int64{
				"Assasinu Credi": 19,
				"Dono da Bola":   14,
				"Isgalamido":     17,
				"Mal":            6,
				"Oootsimo":       22,
				"Zeh":            19,
			},
			KillsByMeans: map[string]uint{
				"MOD_FALLING":       3,
				"MOD_MACHINEGUN":    4,
				"MOD_RAILGUN":       9,
				"MOD_ROCKET":        37,
				"MOD_ROCKET_SPLASH": 60,
				"MOD_SHOTGUN":       4,
				"MOD_TRIGGER_HURT":  14,
			},
		},
	},
}

func TestExec(t *testing.T) {
	svcResult := services.QuakeLogParserService{}.Exec("../../externals/quake_3_logs/qgames.log")

	for key, _ := range qlpTest.expeted {
		expected := qlpTest.expeted[key]
		received := svcResult[key]

		if expected.TotalKills != received.TotalKills {
			t.Errorf("%s: TotalKills expected %d, received %d", key, expected.TotalKills, received.TotalKills)
		}

		sort.Strings(expected.Players)
		sort.Strings(received.Players)
		if !reflect.DeepEqual(expected.Players, received.Players) {
			t.Errorf("%s: expected %s, received %s", key, expected.Players, received.Players)
		}

		if len(expected.Kills) != len(received.Kills) {
			t.Errorf("%s: Number of keys in kills expected %d, received %d", key, len(expected.Kills), len(received.Kills))
		}

		if len(expected.KillsByMeans) != len(received.KillsByMeans) {
			t.Errorf("%s: Number of keys in kills by means expected %d, received %d", key, len(expected.KillsByMeans), len(received.KillsByMeans))
		}

		for killsKeys, _ := range expected.Kills {
			if expected.Kills[killsKeys] != received.Kills[killsKeys] {
				t.Errorf("%s: %s expected %d, received %d", key, killsKeys, expected.Kills[killsKeys], received.Kills[killsKeys])
			}

		}

		for killsByMeans, _ := range expected.KillsByMeans {
			if expected.KillsByMeans[killsByMeans] != received.KillsByMeans[killsByMeans] {
				t.Errorf("%s: %s expected %d, received %d", key, killsByMeans, expected.KillsByMeans[killsByMeans], received.KillsByMeans[killsByMeans])
			}
		}
	}
}
