package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ygortgaleno/cloudwalk-test/internals/entities"
	"github.com/ygortgaleno/cloudwalk-test/pkg"
)

type Killer struct {
	Id uint32
	entities.Player
}

type Victim struct {
	Id uint32
	entities.Player
}

func main() {
	f, err := os.Open("../qgames.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	rgxNewGameEvent := regexp.MustCompile("InitGame: ")

	var s string
	var games []*pkg.AVLTree[entities.Player]
	var currentGame *pkg.AVLTree[entities.Player]
	for scanner.Scan() {
		s = scanner.Text()

		if rgxNewGameEvent.Match([]byte(s)) {
			currentGame = &pkg.AVLTree[entities.Player]{}
			games = append(games, currentGame)
		}

		InsertKillEventsInfoIntoAvl(s, currentGame)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var gamesDto []*GameDto
	var jsons []string
	for i := 0; i < len(games); i++ {
		gamesDto = append(gamesDto, &GameDto{0, []string{}, map[string]int64{}})
		makeGameDto(games[i].Root, gamesDto[i])
		y := map[string]GameDto{fmt.Sprintf("game_%d", i): *gamesDto[i]}
		x, _ := json.Marshal(y)
		jsons = append(jsons, string(x))
	}
	err = os.WriteFile("../n.json", []byte(fmt.Sprintf("[%s]", strings.Join(jsons, ",\n"))), os.ModeAppend)
	fmt.Print(err)

}

func InsertKillEventsInfoIntoAvl(fileline string, currentGame *pkg.AVLTree[entities.Player]) {
	rgxKillEvent := regexp.MustCompile(`Kill: (\d+) (\d+) (\d+):\s(.*) killed (.*) by (\w+)`)
	rgxResult := rgxKillEvent.FindStringSubmatch(fileline)

	if rgxResult != nil {
		idKiller, _ := strconv.ParseUint(rgxResult[1], 10, 64)
		idVictim, _ := strconv.ParseUint(rgxResult[2], 10, 64)
		killer := entities.Player{Name: rgxResult[4], Kills: 1, WorldDeaths: 0}
		victim := entities.Player{Name: rgxResult[5], Kills: 0, WorldDeaths: 0}

		if killer.Name == "<world>" {
			victim.WorldDeaths = 1
		}

		if node := currentGame.Search(uint32(idKiller)); node != nil {
			node.Data.Kills += 1
		} else {
			currentGame.Insert(uint32(idKiller), killer)
		}

		if node := currentGame.Search(uint32(idVictim)); node != nil {
			node.Data.WorldDeaths += victim.WorldDeaths
		} else {
			currentGame.Insert(uint32(idVictim), victim)
		}
	}
}

func makeGameDto(node *pkg.Node[entities.Player], gDto *GameDto) {
	if node == nil {
		return
	}

	if node.Data.Name != "<world>" {
		gDto.Kills[node.Data.Name] = int64(node.Data.Kills) - node.Data.WorldDeaths
		gDto.Players = append(gDto.Players, node.Data.Name)
	}
	gDto.TotalKills += node.Data.Kills

	makeGameDto(node.Left, gDto)
	makeGameDto(node.Right, gDto)
}

type GameDto struct {
	TotalKills uint32           `json:"total_kills"`
	Players    []string         `json:"players"`
	Kills      map[string]int64 `json:"kills"`
}
