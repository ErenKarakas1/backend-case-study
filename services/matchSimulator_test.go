package services

import (
	"math/rand"
	"testing"

	"insider/models"

	"github.com/stretchr/testify/assert"
)

func TestRandomizedMatchSimulator_SimulateMatch(t *testing.T) {
	sim := NewMatchSimulator().(*RandomizedMatchSimulator)
	// deterministic
	sim.random = rand.New(rand.NewSource(42))

	team := models.Team{
		ID:   1,
		Name: "X",
		Attributes: models.TeamAttributes{
			Attack:    0.5,
			Defense:   0.5,
			Midfield:  0.5,
			HomeBoost: 1.0,
		},
		PlayStyle: models.PlayStyleBalanced,
	}

	res := sim.SimulateMatch(team, team)

	assert.GreaterOrEqual(t, res.HomeScore, 0, "home score should be non-negative")
	assert.GreaterOrEqual(t, res.AwayScore, 0, "away score should be non-negative")
	assert.LessOrEqual(t, res.HomeScore, 8, "home score should be <= 8")
	assert.LessOrEqual(t, res.AwayScore, 8, "away score should be <= 8")
}
