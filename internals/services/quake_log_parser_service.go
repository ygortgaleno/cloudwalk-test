package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/ygortgaleno/cloudwalk-test/internals/dtos"
	"github.com/ygortgaleno/cloudwalk-test/internals/entities"
	"github.com/ygortgaleno/cloudwalk-test/internals/pkg"
)

type PlayersAVLTree = pkg.AVLTree[entities.Player]

type DeathCauseAVLTree = pkg.AVLTree[entities.DeathCause]

type QuakeLogParserService struct{}

func (svc QuakeLogParserService) Exec(filepath string) []byte {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rgxNewGameEvent := regexp.MustCompile("InitGame: ")

	var fLine string
	var pTree *PlayersAVLTree
	var dcTree *DeathCauseAVLTree
	var pTrees []*PlayersAVLTree
	var dcTrees []*DeathCauseAVLTree
	var wg sync.WaitGroup

	for scanner.Scan() {
		fLine = scanner.Text()
		if rgxNewGameEvent.Match([]byte(fLine)) {
			pTree = &PlayersAVLTree{}
			dcTree = &DeathCauseAVLTree{}
			pTrees = append(pTrees, pTree)
			dcTrees = append(dcTrees, dcTree)
		}

		wg.Add(1)
		go func(str string, pTree *PlayersAVLTree, dcTree *DeathCauseAVLTree) {
			defer wg.Done()
			svc.extractInfosFromKillEvent(str, pTree, dcTree)
		}(fLine, pTree, dcTree)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	ch := make(chan string)
	wg.Add(len(pTrees))
	for i := 0; i < len(pTrees); i++ {
		go func(gameCounter int, pTree *PlayersAVLTree, dcTree *DeathCauseAVLTree) {
			ch <- svc.transformTreeIntoGameDto(gameCounter, pTree, dcTree)
			defer wg.Done()
		}(i, pTrees[i], dcTrees[i])
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var jsons []string
	for i := range ch {
		jsons = append(jsons, i)
	}

	pTrees[20].Print()

	return []byte(fmt.Sprintf("[%s]", strings.Join(jsons, ",\n")))
}

func (svc QuakeLogParserService) extractInfosFromKillEvent(str string, pTree *PlayersAVLTree, dcTree *DeathCauseAVLTree) {
	rgxKillEvent := regexp.MustCompile(`Kill: (\d+) (\d+) (\d+):\s(.*) killed (.*) by (\w+)`)
	rgxResult := rgxKillEvent.FindStringSubmatch(str)

	if rgxResult != nil {
		idKiller, idVictim, idDeathCause, killerName, victimName, dcName := rgxResult[1], rgxResult[2],
			rgxResult[3], rgxResult[4], rgxResult[5], rgxResult[6]

		var wg sync.WaitGroup
		defer wg.Wait()
		wg.Add(3)

		// Killer insertion goroutine
		go func(idStr string, name string) {
			defer wg.Done()
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			if node := pTree.Search(uint32(id)); node != nil {
				node.Data.UpdatePlayerKill()
			} else {
				pTree.Insert(uint32(id), entities.Player{Name: name, Kills: 1, WorldDeaths: 0})
			}
		}(idKiller, killerName)

		// Victim insertion goroutine
		go func(idStr string, name string, killerName string) {
			defer wg.Done()

			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			worldDeath := 0
			if killerName == "<world>" {
				worldDeath++
			}

			if node := pTree.Search(uint32(id)); node != nil {
				node.Data.UpdateWorldDeaths(int64(worldDeath))
			} else {
				pTree.Insert(uint32(id), entities.Player{Name: name, Kills: 0, WorldDeaths: int64(worldDeath)})
			}
		}(idVictim, victimName, killerName)

		// DeathCause insertion goroutine
		go func(idStr string, name string) {
			defer wg.Done()

			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			if node := dcTree.Search(uint32(id)); node != nil {
				node.Data.UpdateCounter()
			} else {
				dcTree.Insert(uint32(id), entities.DeathCause{Name: name, Counter: 1})
			}
		}(idDeathCause, dcName)
	}
}

func (svc QuakeLogParserService) transformTreeIntoGameDto(gameCounter int, pTree *PlayersAVLTree, dcTree *DeathCauseAVLTree) string {
	gameDto := &dtos.GameDto{TotalKills: 0, Players: []string{}, Kills: map[string]int64{}, KillsByMeans: map[string]uint{}}
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		svc.putPlayersTreeInfoIntoGameDto(pTree.Root, gameDto)
	}()

	go func() {
		defer wg.Done()
		svc.putDeathCausesTreeInfoIntoGameDto(dcTree.Root, gameDto)
	}()
	wg.Wait()
	unmarsheledJson := map[string]dtos.GameDto{fmt.Sprintf("game_%d", gameCounter): *gameDto}
	jsonBytes, _ := json.Marshal(unmarsheledJson)
	return string(jsonBytes)
}

func (svc QuakeLogParserService) putPlayersTreeInfoIntoGameDto(pNode *pkg.Node[entities.Player], gDto *dtos.GameDto) {
	if pNode == nil {
		return
	}

	if pNode.Data.Name != "<world>" {
		gDto.Kills[pNode.Data.Name] = int64(pNode.Data.Kills) - pNode.Data.WorldDeaths
		gDto.Players = append(gDto.Players, pNode.Data.Name)
	}
	gDto.TotalKills += pNode.Data.Kills

	svc.putPlayersTreeInfoIntoGameDto(pNode.Left, gDto)
	svc.putPlayersTreeInfoIntoGameDto(pNode.Right, gDto)
}

func (svc QuakeLogParserService) putDeathCausesTreeInfoIntoGameDto(dcNode *pkg.Node[entities.DeathCause], gDto *dtos.GameDto) {
	if dcNode == nil {
		return
	}

	gDto.KillsByMeans[dcNode.Data.Name] = dcNode.Data.Counter

	svc.putDeathCausesTreeInfoIntoGameDto(dcNode.Left, gDto)
	svc.putDeathCausesTreeInfoIntoGameDto(dcNode.Right, gDto)
}
