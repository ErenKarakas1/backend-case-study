package services

import (
	"sort"

	"insider/models"
)

type DefaultLeagueTable struct {
	entryMap map[int]*models.LeagueTableEntry
}

func NewLeagueTable(teams []models.Team) *DefaultLeagueTable {
	entries := make(map[int]*models.LeagueTableEntry)
	for _, team := range teams {
		entries[team.ID] = &models.LeagueTableEntry{
			Team:         team,
			Position:     0,
			Played:       0,
			Won:          0,
			Drawn:        0,
			Lost:         0,
			GoalsFor:     0,
			GoalsAgainst: 0,
			GoalDiff:     0,
			Points:       0,
		}
	}

	return &DefaultLeagueTable{
		entryMap: entries,
	}
}

func (lt *DefaultLeagueTable) CalculateTable(matches []models.Match) []models.LeagueTableEntry {
	for _, e := range lt.entryMap {
		e.Position = 0
		e.Played = 0
		e.Won = 0
		e.Drawn = 0
		e.Lost = 0
		e.GoalsFor = 0
		e.GoalsAgainst = 0
		e.GoalDiff = 0
		e.Points = 0
	}

	for _, match := range matches {
		if !match.IsPlayed {
			continue
		}

		homeEntry := lt.entryMap[match.HomeTeam.ID]
		awayEntry := lt.entryMap[match.AwayTeam.ID]

		if homeEntry == nil || awayEntry == nil {
			continue // Skip if either team entry is missing
		}

		homeEntry.Played++
		awayEntry.Played++

		result := match.Result
		hs, as := result.HomeScore, result.AwayScore

		homeEntry.GoalsFor += hs
		homeEntry.GoalsAgainst += as
		awayEntry.GoalsFor += as
		awayEntry.GoalsAgainst += hs

		if result.IsWin() {
			homeEntry.Won++
			awayEntry.Lost++
			homeEntry.Points += 3
		} else if result.IsDraw() {
			homeEntry.Drawn++
			awayEntry.Drawn++
			homeEntry.Points += 1
			awayEntry.Points += 1
		} else {
			awayEntry.Won++
			homeEntry.Lost++
			awayEntry.Points += 3
		}
	}

	var table []models.LeagueTableEntry
	for _, entry := range lt.entryMap {
		entry.GoalDiff = entry.GoalsFor - entry.GoalsAgainst
		table = append(table, *entry)
	}

	// Sort entries by points, then goal difference, then goals for
	sort.Slice(table, func(i, j int) bool {
		a, b := table[i], table[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GoalDiff != b.GoalDiff {
			return a.GoalDiff > b.GoalDiff
		}
		return a.GoalsFor > b.GoalsFor
	})

	for i := range table {
		table[i].Position = i + 1
	}

	return table
}
