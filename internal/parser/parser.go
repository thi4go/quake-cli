package parser

import (
	"bufio"
	"os"
	"regexp"

	"github.com/thi4go/quake-cli/internal/models"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseLogFile(filePath string) ([]models.Game, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var games []models.Game
	var currentGame *models.Game
	scanner := bufio.NewScanner(file)

	initGameRegex := regexp.MustCompile(`InitGame:`)
	killRegex := regexp.MustCompile(`Kill: (\d+) (\d+) (\d+): (.+) killed (.+) by (.+)`)
	clientUserinfoChangedRegex := regexp.MustCompile(`ClientUserinfoChanged: (\d+) n\\(.+)\\t\\`)

	for scanner.Scan() {
		line := scanner.Text()

		if initGameRegex.MatchString(line) {
			if currentGame != nil {
				games = append(games, *currentGame)
			}
			currentGame = models.NewGame(len(games) + 1)
		}

		if currentGame == nil {
			continue
		}

		if killMatch := killRegex.FindStringSubmatch(line); killMatch != nil {
			killer := killMatch[4]
			victim := killMatch[5]
			mod := killMatch[6]

			currentGame.TotalKills++
			currentGame.KillsByMeans[mod]++

			if killer != "<world>" && killer != victim {
				currentGame.AddKill(killer)
			}

			currentGame.AddPlayer(victim)
			if killer != "<world>" {
				currentGame.AddPlayer(killer)
			}
		}

		if clientMatch := clientUserinfoChangedRegex.FindStringSubmatch(line); clientMatch != nil {
			playerName := clientMatch[2]
			currentGame.AddPlayer(playerName)
		}
	}

	if currentGame != nil {
		games = append(games, *currentGame)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return games, nil
}
