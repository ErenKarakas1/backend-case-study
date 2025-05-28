package services

import (
	"math/rand"
	"time"

	"insider/models"
)

type RoundRobinScheduler struct {
	random *rand.Rand
}

func NewMatchScheduler() MatchScheduler {
	return &RoundRobinScheduler{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateSchedule creates a round-robin schedule for the given teams.
// Assumes an even number of teams since task specifies 4 teams anyway.
func (s *RoundRobinScheduler) GenerateSchedule(teams []models.Team) []models.Match {
	n := len(teams)

	s.random.Shuffle(n, func(i, j int) {
		teams[i], teams[j] = teams[j], teams[i]
	})

	matches := make([]models.Match, 0, n*(n-1))

	matchID := 1

	weekCount := n * (n - 1) / 2
	for week := range weekCount {
		for i := range 2 {
			home := teams[i]
			away := teams[n-1-i]

			if week%2 == 1 {
				home, away = away, home
			}

			match := models.Match{
				ID:       matchID,
				Week:     week + 1,
				HomeTeam: &home,
				AwayTeam: &away,
				IsPlayed: false,
			}

			matchID++
			matches = append(matches, match)
		}

		// Rotate teams for next week
		teams = append(teams[:1], append([]models.Team{teams[n-1]}, teams[1:n-1]...)...)
	}

	return matches
}
