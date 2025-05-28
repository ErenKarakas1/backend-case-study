package services

import (
	"testing"

	"insider/models"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLeagueTable_CalculateTable(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	}

	lt := NewLeagueTable(teams)
	matches := []models.Match{
		{
			HomeTeam: &teams[0],
			AwayTeam: &teams[1],
			IsPlayed: true,
			Result:   models.MatchResult{HomeScore: 2, AwayScore: 1},
		},
		{
			HomeTeam: &teams[1],
			AwayTeam: &teams[0],
			IsPlayed: true,
			Result:   models.MatchResult{HomeScore: 1, AwayScore: 1},
		},
	}

	table := lt.CalculateTable(matches)
	assert.Len(t, table, 2, "expected 2 teams in the table")

	var a, b models.LeagueTableEntry
	for _, e := range table {
		if e.Team.ID == 1 {
			a = e
		}
		if e.Team.ID == 2 {
			b = e
		}
	}

	// Team A: 1W 1D -> 4 pts, GF=3 GA=2
	assert.Equal(t, 4, a.Points, "Team A should have 4 points")
	assert.Equal(t, 3, a.GoalsFor, "Team A should have 3 goals for")
	assert.Equal(t, 2, a.GoalsAgainst, "Team A should have 2 goals against")

	// Team B: 1D 1L -> 1 pt, GF=2 GA=3
	assert.Equal(t, 1, b.Points, "Team B should have 1 point")
	assert.Equal(t, 2, b.GoalsFor, "Team B should have 2 goals for")
	assert.Equal(t, 3, b.GoalsAgainst, "Team B should have 3 goals against")
}

func TestDefaultLeagueTable_ResetsBetweenCalls(t *testing.T) {
	teams := []models.Team{{ID: 1}, {ID: 2}}
	lt := NewLeagueTable(teams)

	// first call: one match
	m1 := []models.Match{{
		HomeTeam: &teams[0], AwayTeam: &teams[1], IsPlayed: true,
		Result: models.MatchResult{HomeScore: 3, AwayScore: 1},
	}}
	_ = lt.CalculateTable(m1)

	// second call: empty (simulate reset)
	tbl2 := lt.CalculateTable(nil)
	for _, e := range tbl2 {
		assert.Equal(t, 0, e.Played, "expected Played to be reset to 0")
		assert.Equal(t, 0, e.Points, "expected Points to be reset to 0")
		assert.Equal(t, 0, e.GoalsFor, "expected GoalsFor to be reset to 0")
	}

	// third call: back to one match
	tbl3 := lt.CalculateTable(m1)
	var a3 models.LeagueTableEntry
	for _, e := range tbl3 {
		if e.Team.ID == 1 {
			a3 = e
		}
	}

	assert.Equal(t, 1, a3.Played, "expected Played to be 1 after third call")
	assert.Equal(t, 3, a3.Points, "expected Points to be 3 after third call")
	assert.Equal(t, 3, a3.GoalsFor, "expected GoalsFor to be 3 after third call")
}
