package models

type Game struct {
	ID           int
	TotalKills   int
	Players      map[string]bool
	Kills        map[string]int
	KillsByMeans map[string]int
	KillSequence []Kill
}

type Kill struct {
	Killer string
	Victim string
	Weapon string
}

func NewGame(id int) *Game {
	return &Game{
		ID:           id,
		Players:      make(map[string]bool),
		Kills:        make(map[string]int),
		KillsByMeans: make(map[string]int),
	}
}

func (g *Game) AddKill(player string) {
	g.Kills[player]++
}

func (g *Game) AddPlayer(name string) {
	if _, exists := g.Players[name]; !exists {
		g.Players[name] = true
		g.Kills[name] = 0
	}
}
func (g *Game) AddKillToSequence(killer, victim, weapon string) {
	g.KillSequence = append(g.KillSequence, Kill{
		Killer: killer,
		Victim: victim,
		Weapon: weapon,
	})
}
