package services

import (
	"testing"

	"insider/models"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinScheduler_GenerateSchedule(t *testing.T) {
	teams := []models.Team{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
		{ID: 4, Name: "D"},
	}
	scheduler := NewMatchScheduler()
	matches := scheduler.GenerateSchedule(teams)

	expected := (len(teams) - 1) * len(teams)
	assert.Equal(t, expected, len(matches), "expected %d matches, got %d", expected, len(matches))

	seen := make(map[int]struct{})
	for _, m := range matches {
		assert.NotEqual(t, m.HomeTeam.ID, m.AwayTeam.ID, "home and away teams must differ")
		assert.Greater(t, m.HomeTeam.ID, 0, "home team ID must be positive")

		assert.Greater(t, m.Week, 0, "week must be positive")
		assert.LessOrEqual(t, m.Week, 2*expected/len(teams), "week %d should be less than %d", m.Week, len(teams))

		key := m.Week<<16 | m.HomeTeam.ID<<8 | m.AwayTeam.ID
		assert.NotContains(t, seen, key, "duplicate match found: %#v", m)
		seen[key] = struct{}{}
	}
}
