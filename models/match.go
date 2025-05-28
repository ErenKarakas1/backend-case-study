package models

type MatchResult struct {
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}

type Match struct {
	ID       int         `json:"id"`
	Week     int         `json:"week"`
	HomeTeam *Team       `json:"home_team"`
	AwayTeam *Team       `json:"away_team"`
	Result   MatchResult `json:"result"`
	IsPlayed bool        `json:"is_played"`
}

func (mr MatchResult) IsWin() bool {
	return mr.HomeScore > mr.AwayScore
}

func (mr MatchResult) IsDraw() bool {
	return mr.HomeScore == mr.AwayScore
}

func (mr MatchResult) IsLoss() bool {
	return mr.HomeScore < mr.AwayScore
}
