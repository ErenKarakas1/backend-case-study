package services

import "insider/models"

// MatchSimulator defines the interface for simulating match results
type MatchSimulator interface {
	SimulateMatch(homeTeam, awayTeam models.Team) models.MatchResult
}

// LeagueTable defines the interface for calculating league tables
type LeagueTable interface {
	CalculateTable(matches []models.Match) []models.LeagueTableEntry
}

// LeaguePredictor defines the interface for predicting the championship odds for teams
type LeaguePredictor interface {
	CalculateChampionshipOdds(table []models.LeagueTableEntry, remaining []models.Match) []models.ChampionshipOdds
}

// MatchScheduler defines the interface for generating match schedules
type MatchScheduler interface {
	GenerateSchedule(teams []models.Team) []models.Match
}

// LeagueService defines the main service interface
type LeagueService interface {
	GetCurrentState() (*models.LeagueSimulation, error)
	SimulateNextWeek() (*models.WeekSimulation, error)
	SimulateRemainingWeeks() (*models.LeagueSimulation, error)
	ResetSimulation() error
	UpdateMatchResult(matchID int, homeScore, awayScore int) error
}
