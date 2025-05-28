package services

import (
	"math"
	"math/rand"
	"time"

	"insider/models"
)

type matchupMatrix struct {
	advantages [][]float64
	styleIndex map[models.PlayStyle]int
}

type RandomizedMatchSimulator struct {
	random        *rand.Rand
	matchupMatrix matchupMatrix
}

func NewMatchSimulator() MatchSimulator {
	return &RandomizedMatchSimulator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		matchupMatrix: matchupMatrix{
			// extra zeroes at the end are there for padding
			advantages: [][]float64{
				//   	   A      D      P     B
				/* A */ {0.000, -0.15, 0.050, 0.10},
				/* D */ {0.150, 0.000, -0.10, 0.05},
				/* P */ {-0.05, 0.100, 0.000, 0.10},
				/* B */ {-0.05, -0.05, -0.05, 0.00},
			},
			styleIndex: map[models.PlayStyle]int{
				models.PlayStyleAttacking:  0,
				models.PlayStyleDefensive:  1,
				models.PlayStylePossession: 2,
				models.PlayStyleBalanced:   3,
			},
		},
	}
}

func (sim *RandomizedMatchSimulator) SimulateMatch(home, away models.Team) models.MatchResult {
	headToHead := sim.calculateHeadToHeadMatchup(home, away)

	homeExpectedGoals := sim.calculateExpectedGoals(home, away, true, headToHead)
	awayExpectedGoals := sim.calculateExpectedGoals(away, home, false, -headToHead)

	// Add some variance to the expected goals
	homeGoals := sim.simulateGoalsFromExpected(homeExpectedGoals)
	awayGoals := sim.simulateGoalsFromExpected(awayExpectedGoals)

	return models.MatchResult{
		HomeScore: homeGoals,
		AwayScore: awayGoals,
	}
}

func (sim *RandomizedMatchSimulator) calculateHeadToHeadMatchup(home, away models.Team) float64 {
	row, ok1 := sim.matchupMatrix.styleIndex[home.PlayStyle]
	col, ok2 := sim.matchupMatrix.styleIndex[away.PlayStyle]
	if ok1 && ok2 {
		return sim.matchupMatrix.advantages[row][col]
	}

	return 0.0 // No advantage, should not happen
}

func (sim *RandomizedMatchSimulator) calculateExpectedGoals(attackingTeam, defendingTeam models.Team, isHome bool, headToHead float64) float64 {
	attack := attackingTeam.Attributes.Attack
	midfield := attackingTeam.Attributes.Midfield * 0.5

	homeAdvantage := 0.0
	if isHome {
		homeAdvantage = attackingTeam.Attributes.HomeBoost
	}

	resistance := defendingTeam.Attributes.Defense

	expected := 1.0 + (attack + midfield - resistance)

	expected += headToHead
	expected += homeAdvantage

	if expected < 0.1 {
		expected = 0.1
	}
	if expected > 5.0 {
		expected = 5.0
	}

	return expected
}

func (sim *RandomizedMatchSimulator) simulateGoalsFromExpected(expectedGoals float64) int {
	goals := 0
	L := math.Exp(-expectedGoals)
	p := 1.0

	for p > L && goals < 8 { // Limit to 8 goals
		goals++
		p *= sim.random.Float64()
	}

	return goals - 1
}
