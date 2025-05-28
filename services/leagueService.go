package services

import (
	"log"

	"insider/database"
	"insider/models"
)

type BasicLeagueService struct {
	db             database.Database
	matchSimulator MatchSimulator
	table          LeagueTable
	matchScheduler MatchScheduler
	predictor      LeaguePredictor
	teamMap        map[int]models.Team
}

func NewLeagueService(db database.Database, simulator MatchSimulator, table LeagueTable, scheduler MatchScheduler, predictor LeaguePredictor) LeagueService {
	teams, err := db.GetTeams()
	if err != nil {
		panic("Failed to retrieve teams from database: " + err.Error())
	}

	teamMap := make(map[int]models.Team)
	for _, team := range teams {
		teamMap[team.ID] = team
	}

	return &BasicLeagueService{
		db:             db,
		matchSimulator: simulator,
		table:          table,
		matchScheduler: scheduler,
		predictor:      predictor,
		teamMap:        teamMap,
	}
}

func (ls *BasicLeagueService) GetCurrentState() (*models.LeagueSimulation, error) {
	state, err := ls.db.GetSimulationState()
	if err != nil {
		return nil, err
	}

	matches, err := ls.db.GetMatches()
	if err != nil {
		return nil, err
	}

	table := ls.table.CalculateTable(matches)

	simulation := &models.LeagueSimulation{
		CurrentWeek: state.CurrentWeek,
		MaxWeeks:    state.MaxWeeks,
		Table:       table,
		Matches:     matches,
	}

	if state.CurrentWeek > 4 {
		remainingMatches := ls.getRemainingMatches(matches)
		if len(remainingMatches) > 0 {
			odds := ls.predictor.CalculateChampionshipOdds(table, remainingMatches)
			simulation.ChampionshipOdds = odds
		}
	}
	return simulation, nil
}

func (ls *BasicLeagueService) SimulateNextWeek() (*models.WeekSimulation, error) {
	state, err := ls.db.GetSimulationState()
	if err != nil {
		return nil, err
	}

	if state.CurrentWeek > state.MaxWeeks {
		sim, err := ls.GetCurrentState()
		if err != nil {
			return nil, err
		}
		return ls.filterStateByWeek(sim, sim.CurrentWeek-1), nil
	}

	weekMatches, err := ls.db.GetMatchesForWeek(state.CurrentWeek)
	if err != nil {
		return nil, err
	}

	for _, match := range weekMatches {
		if !match.IsPlayed {
			homeTeam := ls.teamMap[match.HomeTeam.ID]
			awayTeam := ls.teamMap[match.AwayTeam.ID]

			result := ls.matchSimulator.SimulateMatch(homeTeam, awayTeam)

			err = ls.db.UpdateMatchResult(match.ID, result)
			if err != nil {
				return nil, err
			}
		}
	}

	err = ls.db.UpdateCurrentWeek(state.CurrentWeek + 1)
	if err != nil {
		return nil, err
	}

	sim, err := ls.GetCurrentState()
	if err != nil {
		return nil, err
	}
	return ls.filterStateByWeek(sim, sim.CurrentWeek-1), nil
}

func (ls *BasicLeagueService) SimulateRemainingWeeks() (*models.LeagueSimulation, error) {
	state, err := ls.db.GetSimulationState()
	if err != nil {
		return nil, err
	}

	if state.CurrentWeek > state.MaxWeeks {
		log.Println("Simulation already completed")
		return ls.GetCurrentState()
	}

	diff := state.MaxWeeks - state.CurrentWeek + 1

	for range diff {
		_, err = ls.SimulateNextWeek()
		if err != nil {
			return nil, err
		}
	}
	return ls.GetCurrentState()
}

func (ls *BasicLeagueService) ResetSimulation() error {
	err := ls.db.ResetSimulation()
	if err != nil {
		return err
	}

	// Regenerate the schedule
	teams, err := ls.db.GetTeams()
	if err != nil {
		return err
	}

	matches := ls.matchScheduler.GenerateSchedule(teams)

	err = ls.db.InsertMatches(matches)
	if err != nil {
		return err
	}
	return nil
}

func (ls *BasicLeagueService) UpdateMatchResult(matchID int, homeScore, awayScore int) error {
	result := models.MatchResult{
		HomeScore: homeScore,
		AwayScore: awayScore,
	}

	if err := ls.db.UpdateMatchResult(matchID, result); err != nil {
		return err
	}
	return nil
}

func (ls *BasicLeagueService) getRemainingMatches(matches []models.Match) []models.Match {
	remaining := make([]models.Match, 0)
	for _, match := range matches {
		if !match.IsPlayed {
			remaining = append(remaining, match)
		}
	}
	return remaining
}

func (ls *BasicLeagueService) filterStateByWeek(state *models.LeagueSimulation, week int) *models.WeekSimulation {
	filteredMatches := make([]models.Match, 0)
	for _, match := range state.Matches {
		if match.Week == week {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return &models.WeekSimulation{
		PlayedWeek: week,
		Matches:    filteredMatches,
	}
}
