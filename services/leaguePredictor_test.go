package services

import (
	"testing"

	"insider/models"

	"github.com/stretchr/testify/assert"
)

type noopSim struct{}

func (noopSim) SimulateMatch(home, away models.Team) models.MatchResult {
	panic("should not be called when no remaining matches")
}

func TestRandomizedPredictor_NoRemainingMatches(t *testing.T) {
	// setup a table where A leads B
	e1 := models.LeagueTableEntry{Team: models.Team{ID: 1, Name: "A"}, Points: 10}
	e2 := models.LeagueTableEntry{Team: models.Team{ID: 2, Name: "B"}, Points: 5}
	table := []models.LeagueTableEntry{e1, e2}

	pred := &RandomizedPredictor{
		simulator: noopSim{},
		table:     nil,
		random:    nil,
	}

	odds := pred.CalculateChampionshipOdds(table, nil)

	assert.Equal(t, 2, len(odds), "should return odds for both teams")

	for _, o := range odds {
		switch o.TeamID {
		case 1:
			assert.Equal(t, 1.0, o.Probability, "team A should have 100%% probability")
		case 2:
			assert.Equal(t, 0.0, o.Probability, "team B should have 0%% probability")
		default:
			assert.Fail(t, "unexpected team ID in odds: %d", o.TeamID)
		}
	}
}
