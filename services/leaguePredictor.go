package services

import (
	"math/rand"
	"time"

	"insider/models"
)

type RandomizedPredictor struct {
	simulator MatchSimulator
	table     LeagueTable
	random    *rand.Rand
}

func NewLeaguePredictor(simulator MatchSimulator, table LeagueTable) LeaguePredictor {
	return &RandomizedPredictor{
		simulator: simulator,
		table:     table,
		random:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (p *RandomizedPredictor) CalculateChampionshipOdds(table []models.LeagueTableEntry, remaining []models.Match) []models.ChampionshipOdds {
	wins := make(map[int]int, len(table))
	const iters int = 10000 // number of simulations, arbitrary

	simTable := make(map[int]*models.LeagueTableEntry, len(table))

	for range iters {
		clear(simTable)

		for _, e := range table {
			copy := e
			simTable[e.Team.ID] = &copy
		}

		for _, m := range remaining {
			homeTeam := simTable[m.HomeTeam.ID]
			awayTeam := simTable[m.AwayTeam.ID]

			result := p.simulator.SimulateMatch(homeTeam.Team, awayTeam.Team)
			hs, as := result.HomeScore, result.AwayScore

			homeTeam.Played++
			awayTeam.Played++
			homeTeam.GoalsFor += hs
			homeTeam.GoalsAgainst += as
			awayTeam.GoalsFor += as
			awayTeam.GoalsAgainst += hs

			if result.IsWin() {
				homeTeam.Won++
				awayTeam.Lost++
				homeTeam.Points += 3
			} else if result.IsDraw() {
				homeTeam.Drawn++
				awayTeam.Drawn++
				homeTeam.Points += 1
				awayTeam.Points += 1
			} else {
				awayTeam.Won++
				homeTeam.Lost++
				awayTeam.Points += 3
			}
		}

		champID := 0
		bestPts := -1
		bestGD := -1_000_000
		bestGF := -1
		for _, e := range simTable {
			gd := e.GoalsFor - e.GoalsAgainst
			if e.Points > bestPts || (e.Points == bestPts && gd > bestGD) ||
				(e.Points == bestPts && gd == bestGD && e.GoalsFor > bestGF) {
				bestPts = e.Points
				bestGD = gd
				bestGF = e.GoalsFor
				champID = e.Team.ID
			}
		}
		wins[champID]++
	}

	out := make([]models.ChampionshipOdds, 0, len(table))
	for _, e := range table {
		prob := float64(wins[e.Team.ID]) / float64(iters)
		out = append(out, models.ChampionshipOdds{
			TeamID:      e.Team.ID,
			TeamName:    e.Team.Name,
			Probability: prob,
		})
	}
	return out
}
