package simulator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/thi4go/quake-cli/internal/models"
)

const (
	mapWidth  = 60
	mapHeight = 30
)

type Player struct {
	Name  string
	X     int
	Y     int
	Kills int
}

type Projectile struct {
	X, Y     int
	DX, DY   int
	Symbol   rune
	Lifespan int
}

type Simulator struct {
	Game        *models.Game
	Players     map[string]*Player
	Map         [][]rune
	Projectiles []*Projectile
}

func NewSimulator(game *models.Game) *Simulator {
	s := &Simulator{
		Game:        game,
		Players:     make(map[string]*Player),
		Map:         make([][]rune, mapHeight),
		Projectiles: make([]*Projectile, 0),
	}

	for i := range s.Map {
		s.Map[i] = make([]rune, mapWidth)
		for j := range s.Map[i] {
			s.Map[i][j] = '.'
		}
	}

	for player := range game.Players {
		s.Players[player] = &Player{
			Name:  player,
			X:     rand.Intn(mapWidth),
			Y:     rand.Intn(mapHeight),
			Kills: game.Kills[player],
		}
	}

	return s
}

func (s *Simulator) Run() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < s.Game.TotalKills; i++ {
		s.movePlayersRandomly()
		s.moveProjectiles()
		s.simulateKill()
		s.drawMap()
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Simulation completed")
}

func (s *Simulator) movePlayersRandomly() {
	for _, player := range s.Players {
		dx := rand.Intn(3) - 1
		dy := rand.Intn(3) - 1
		player.X = (player.X + dx + mapWidth) % mapWidth
		player.Y = (player.Y + dy + mapHeight) % mapHeight
	}
}

func (s *Simulator) moveProjectiles() {
	for i := 0; i < len(s.Projectiles); i++ {
		p := s.Projectiles[i]
		p.X = (p.X + p.DX + mapWidth) % mapWidth
		p.Y = (p.Y + p.DY + mapHeight) % mapHeight
		p.Lifespan--
		if p.Lifespan <= 0 {
			s.Projectiles = append(s.Projectiles[:i], s.Projectiles[i+1:]...)
			i--
		}
	}
}

func (s *Simulator) simulateKill() {
	players := make([]string, 0, len(s.Players))
	for name := range s.Players {
		players = append(players, name)
	}

	killer := players[rand.Intn(len(players))]
	victim := players[rand.Intn(len(players))]

	if killer != victim {
		s.Players[killer].Kills++
		s.spawnProjectile(s.Players[killer], s.Players[victim])
	}
}

func (s *Simulator) spawnProjectile(from, to *Player) {
	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx != 0 {
		dx = dx / abs(dx)
	}
	if dy != 0 {
		dy = dy / abs(dy)
	}

	symbols := []rune{'*', '+', 'o', '•'}
	s.Projectiles = append(s.Projectiles, &Projectile{
		X:        from.X,
		Y:        from.Y,
		DX:       dx,
		DY:       dy,
		Symbol:   symbols[rand.Intn(len(symbols))],
		Lifespan: 10,
	})
}

func (s *Simulator) drawMap() {
	fmt.Print("\033[H\033[2J") // Clear the console

	// Draw the map border
	fmt.Println(string(repeat('-', mapWidth+2)))

	for y := 0; y < mapHeight; y++ {
		fmt.Print("|")
		for x := 0; x < mapWidth; x++ {
			char := '.'
			for _, player := range s.Players {
				if player.X == x && player.Y == y {
					char = rune(player.Name[0])
					break
				}
			}
			for _, proj := range s.Projectiles {
				if proj.X == x && proj.Y == y {
					char = proj.Symbol
					break
				}
			}
			fmt.Print(string(char))
		}
		fmt.Println("|")
	}

	fmt.Println(string(repeat('-', mapWidth+2)))

	// Print player information
	fmt.Println("Players:")
	for _, player := range s.Players {
		fmt.Printf("%c - %s (Kills: %d)\n", player.Name[0], player.Name, player.Kills)
	}

	// Print weapon information
	fmt.Println("\nWeapons:")
	fmt.Println("🔫 🚀 ⚡ 💣")
}

func repeat(r rune, n int) []rune {
	b := make([]rune, n)
	for i := range b {
		b[i] = r
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
