package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/thi4go/quake-cli/internal/models"
	"github.com/thi4go/quake-cli/internal/parser"
	"github.com/thi4go/quake-cli/internal/simulator"
)

func main() {
	logFile := flag.String("file", "", "Path to the Quake log file")
	simulate := flag.Bool("simulate", false, "Run simulation")
	flag.Parse()

	if *logFile == "" {
		log.Fatal("Please provide a log file path using the -file flag")
	}

	p := parser.NewParser()
	games, err := p.ParseLogFile(*logFile)
	if err != nil {
		log.Fatalf("Error parsing log file: %v", err)
	}

	sort.Slice(games, func(i, j int) bool {
		return games[i].ID < games[j].ID
	})

	if *simulate {
		runSimulation(games)
	} else {
		generateJSONReport(games)
	}
}

func generateJSONReport(games []models.Game) {
	report := generateReport(games)
	jsonReport, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		log.Fatalf("Error generating JSON report: %v", err)
	}
	fmt.Println(string(jsonReport))
}

func runSimulation(games []models.Game) {
	reader := bufio.NewReader(os.Stdin)
	for i, game := range games {
		fmt.Printf("Simulating game %d\n", game.ID)
		sim := simulator.NewSimulator(&game)
		sim.Run()

		if i < len(games)-1 {
			fmt.Print("Press Enter to continue to the next game...")
			reader.ReadString('\n')
		}
	}
}

func generateReport(games []models.Game) []map[string]interface{} {
	var report []map[string]interface{}

	for _, game := range games {
		gameReport := map[string]interface{}{
			"game_id":        fmt.Sprintf("game_%d", game.ID),
			"total_kills":    game.TotalKills,
			"players":        getPlayersList(game.Players),
			"kills":          game.Kills,
			"kills_by_means": game.KillsByMeans,
		}
		report = append(report, gameReport)
	}

	return report
}

func getPlayersList(players map[string]bool) []string {
	playersList := make([]string, 0, len(players))
	for player := range players {
		playersList = append(playersList, player)
	}
	sort.Strings(playersList)
	return playersList
}
