package services

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/ygortgaleno/cloudwalk-test/internals/dtos"
	"github.com/ygortgaleno/cloudwalk-test/internals/entities"
	"github.com/ygortgaleno/cloudwalk-test/internals/pkg"
)

type PlayersAVLTree = pkg.AVLTree[entities.Player]

type DeathCauseAVLTree = pkg.AVLTree[entities.DeathCause]

type GameInfos struct {
	*PlayersAVLTree
	*DeathCauseAVLTree
}

type GameInfosAVLTree = pkg.AVLTree[GameInfos]

type QuakeLogParserService struct{}

func (svc QuakeLogParserService) Exec(filepath string) map[string]dtos.GameDto {
	file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rgxNewGameEvent := regexp.MustCompile("InitGame: ")
	rgxKillEvent := regexp.MustCompile(`Kill: (\d+) (\d+) (\d+):\s(.*) killed (.*) by (\w+)`)

	var fLine string
	gTree := &GameInfosAVLTree{}
	gameCounter := 0

	for scanner.Scan() {
		fLine = scanner.Text()
		if rgxNewGameEvent.Match([]byte(fLine)) {
			gameCounter++
			gTree.Insert(uint32(gameCounter), GameInfos{&PlayersAVLTree{}, &DeathCauseAVLTree{}})
		} else if killEvent := rgxKillEvent.FindStringSubmatch(fLine); killEvent != nil {
			svc.extractInfosFromKillEvent(killEvent, gTree.Search(uint32(gameCounter)))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	gameMap := make(map[string]dtos.GameDto)
	svc.transformTreeIntoMapGameDto(gTree.Root, &gameMap)

	return gameMap
}

func (svc QuakeLogParserService) extractInfosFromKillEvent(killEvent []string, game *pkg.Node[GameInfos]) {

	if killEvent != nil {
		idKiller, idVictim, idDeathCause, killerName, victimName, dcName := killEvent[1], killEvent[2],
			killEvent[3], killEvent[4], killEvent[5], killEvent[6]

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			defer wg.Done()
			idK, err := strconv.ParseUint(idKiller, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			if node := game.Data.PlayersAVLTree.Search(uint32(idK)); node != nil {
				node.Data.UpdatePlayerKill()
			} else {
				game.Data.PlayersAVLTree.Insert(uint32(idK), entities.Player{Name: killerName, Kills: 1, WorldDeaths: 0})
			}
		}()

		go func() {
			defer wg.Done()
			idV, err := strconv.ParseUint(idVictim, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			worldDeath := 0
			if killerName == "<world>" {
				worldDeath++
			}

			if node := game.Data.PlayersAVLTree.Search(uint32(idV)); node != nil {
				node.Data.UpdateWorldDeaths(int64(worldDeath))
			} else {
				game.Data.PlayersAVLTree.Insert(
					uint32(idV),
					entities.Player{Name: victimName, Kills: 0, WorldDeaths: int64(worldDeath)},
				)
			}
		}()

		go func() {
			defer wg.Done()
			idDc, err := strconv.ParseUint(idDeathCause, 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			if node := game.Data.DeathCauseAVLTree.Search(uint32(idDc)); node != nil {
				node.Data.UpdateCounter()
			} else {
				game.Data.DeathCauseAVLTree.Insert(
					uint32(idDc),
					entities.DeathCause{Name: dcName, Counter: 1},
				)
			}
		}()

		wg.Wait()
	}
}

func (svc QuakeLogParserService) transformTreeIntoMapGameDto(gameInfo *pkg.Node[GameInfos], gamesMap *map[string]dtos.GameDto) {
	if gameInfo == nil {
		return
	}

	svc.transformTreeIntoMapGameDto(gameInfo.Left, gamesMap)

	func() {
		gameDto := &dtos.GameDto{TotalKills: 0, Players: []string{}, Kills: map[string]int64{}, KillsByMeans: map[string]uint{}}
		var wg sync.WaitGroup
		wg.Add(2)

		go func(pTree *PlayersAVLTree, gameDto *dtos.GameDto) {
			defer wg.Done()
			svc.putPlayersTreeInfoIntoGameDto(pTree.Root, gameDto)
		}(gameInfo.Data.PlayersAVLTree, gameDto)

		go func(dcTree *DeathCauseAVLTree, gameDto *dtos.GameDto) {
			defer wg.Done()
			svc.putDeathCausesTreeInfoIntoGameDto(dcTree.Root, gameDto)
		}(gameInfo.Data.DeathCauseAVLTree, gameDto)
		wg.Wait()

		(*gamesMap)[fmt.Sprintf("game_%d", gameInfo.Id)] = *gameDto
	}()

	svc.transformTreeIntoMapGameDto(gameInfo.Right, gamesMap)
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
