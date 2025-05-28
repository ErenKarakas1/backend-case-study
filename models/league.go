package models

type LeagueTableEntry struct {
	Position     int  `json:"position"`
	Team         Team `json:"team"`
	Played       int  `json:"played"`
	Won          int  `json:"won"`
	Drawn        int  `json:"drawn"`
	Lost         int  `json:"lost"`
	GoalsFor     int  `json:"goals_for"`
	GoalsAgainst int  `json:"goals_against"`
	GoalDiff     int  `json:"goal_diff"`
	Points       int  `json:"points"`
}

type ChampionshipOdds struct {
	TeamID      int     `json:"team_id"`
	TeamName    string  `json:"team_name"`
	Probability float64 `json:"probability"`
}

type WeekSimulation struct {
	PlayedWeek int     `json:"played_week"`
	Matches    []Match `json:"matches"`
}

type LeagueSimulation struct {
	CurrentWeek      int                `json:"current_week"`
	MaxWeeks         int                `json:"max_weeks"`
	Table            []LeagueTableEntry `json:"table"`
	Matches          []Match            `json:"matches"`
	ChampionshipOdds []ChampionshipOdds `json:"championship_odds,omitempty"`
}

type SimulationState struct {
	ID          int `json:"id"`
	CurrentWeek int `json:"current_week"`
	MaxWeeks    int `json:"max_weeks"`
}
