package database

import (
	"database/sql"
	"log"

	"insider/models"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Database interface {
	Initialize()
	Close() error

	GetTeams() ([]models.Team, error)
	GetMatches() ([]models.Match, error)
	GetMatchesForWeek(week int) ([]models.Match, error)
	GetSimulationState() (*models.SimulationState, error)

	InsertMatches(matches []models.Match) error

	UpdateMatchResult(matchID int, result models.MatchResult) error
	UpdateCurrentWeek(week int) error

	ResetSimulation() error
}

type SQLiteDatabase struct {
	db   *sql.DB
	path string
}

func NewSQLiteDatabase(databasePath string) *SQLiteDatabase {
	return &SQLiteDatabase{
		db:   nil,
		path: databasePath,
	}
}

func (sqlite *SQLiteDatabase) Initialize() {
	db, err := sql.Open("sqlite3", sqlite.path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	sqlite.db = db

	const createTablesQuery string = `
	CREATE TABLE IF NOT EXISTS teams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		attack REAL NOT NULL DEFAULT 0.5,
		defense REAL NOT NULL DEFAULT 0.5,
		midfield REAL NOT NULL DEFAULT 0.5,
		home_boost REAL NOT NULL DEFAULT 1.0,
		play_style TEXT NOT NULL DEFAULT 'balanced'
	);
	
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		week INTEGER NOT NULL,
		home_team_id INTEGER NOT NULL,
		away_team_id INTEGER NOT NULL,
		home_score INTEGER,
		away_score INTEGER,
		is_played BOOLEAN NOT NULL DEFAULT FALSE,
		FOREIGN KEY (home_team_id) REFERENCES teams(id),
		FOREIGN KEY (away_team_id) REFERENCES teams(id)
	);

	CREATE TABLE IF NOT EXISTS simulation_state (
		id INTEGER PRIMARY KEY DEFAULT 1,
		current_week INTEGER NOT NULL DEFAULT 1,
		max_weeks INTEGER NOT NULL DEFAULT 6
	);
	`
	_, err = sqlite.db.Exec(createTablesQuery)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	const insertTeamsQuery string = `
	INSERT OR IGNORE INTO teams
	(name, attack, defense, midfield, home_boost, play_style) VALUES
	('Manchester City', 0.95, 0.85, 0.92, 0.42, 'possession'),
	('Liverpool', 0.90, 0.80, 0.85, 0.45, 'attacking'),
	('Arsenal', 0.85, 0.75, 0.88, 0.40, 'balanced'),
	('Chelsea', 0.78, 0.88, 0.80, 0.37, 'defensive');
	`
	_, err = sqlite.db.Exec(insertTeamsQuery)
	if err != nil {
		log.Fatalf("Failed to insert initial teams list: %v", err)
	}

	const insertStateQuery string = `
	INSERT OR IGNORE INTO simulation_state (id, current_week, max_weeks)
	VALUES (1, 1, 6);
	`
	_, err = sqlite.db.Exec(insertStateQuery)
	if err != nil {
		log.Fatalf("Failed to insert initial state: %v", err)
	}
}

func (sqlite *SQLiteDatabase) Close() error {
	if sqlite.db != nil {
		err := sqlite.db.Close()
		if err != nil {
			return err
		}
		sqlite.db = nil
	}
	return nil
}

func (sqlite *SQLiteDatabase) GetTeams() ([]models.Team, error) {
	rows, err := sqlite.db.Query(getTeamsQuery)
	if err != nil {
		log.Printf("Failed to query teams: %v", err)
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var team models.Team
		err := rows.Scan(&team.ID, &team.Name, &team.Attributes.Attack, &team.Attributes.Defense,
			&team.Attributes.Midfield, &team.Attributes.HomeBoost, &team.PlayStyle)
		if err != nil {
			log.Printf("Failed to scan team row: %v", err)
			return nil, err
		}
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v", err)
		return nil, err
	}
	return teams, nil
}

func (sqlite *SQLiteDatabase) GetMatches() ([]models.Match, error) {
	rows, err := sqlite.db.Query(getMatchesQuery)

	if err != nil {
		log.Printf("Failed to query matches: %v", err)
		return nil, err
	}
	defer rows.Close()

	return sqlite.populateMatches(rows)
}

func (sqlite *SQLiteDatabase) GetMatchesForWeek(week int) ([]models.Match, error) {
	rows, err := sqlite.db.Query(getMatchesForWeekQuery, week)

	if err != nil {
		log.Printf("Failed to query matches for week %d: %v", week, err)
		return nil, err
	}
	defer rows.Close()

	return sqlite.populateMatches(rows)
}

func (sqlite *SQLiteDatabase) GetSimulationState() (*models.SimulationState, error) {
	var state models.SimulationState
	if err := sqlite.db.QueryRow(getStateQuery).
		Scan(&state.ID, &state.CurrentWeek, &state.MaxWeeks); err != nil {
		log.Printf("Failed to retrieve simulation state: %v", err)
		return nil, err
	}
	return &state, nil
}

func (sqlite *SQLiteDatabase) InsertMatches(matches []models.Match) error {
	if len(matches) == 0 {
		log.Println("No matches to insert")
		return nil
	}

	tx, err := sqlite.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(insertMatchQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, match := range matches {
		if _, err := stmt.Exec(match.Week, match.HomeTeam.ID, match.AwayTeam.ID); err != nil {
			log.Printf("Failed to insert match for week %d between team %d and team %d: %v",
				match.Week, match.HomeTeam.ID, match.AwayTeam.ID, err)
			return err
		}
	}

	return tx.Commit()
}

func (sqlite *SQLiteDatabase) UpdateMatchResult(matchID int, result models.MatchResult) error {
	if _, err := sqlite.db.Exec(updateMatchQuery, result.HomeScore, result.AwayScore, matchID); err != nil {
		log.Printf("Failed to update match result for match ID %d: %v", matchID, err)
		return err
	}
	return nil
}

func (sqlite *SQLiteDatabase) UpdateCurrentWeek(week int) error {
	if _, err := sqlite.db.Exec(updateWeekQuery, week); err != nil {
		log.Printf("Failed to update current week to %d: %v", week, err)
		return err
	}
	return nil
}

func (sqlite *SQLiteDatabase) ResetSimulation() error {
	tx, err := sqlite.db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for reset: %v", err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(resetMatchesQuery)
	if err != nil {
		log.Printf("Failed to reset matches: %v", err)
		return err
	}

	_, err = tx.Exec(resetStateQuery)
	if err != nil {
		log.Printf("Failed to reset simulation state: %v", err)
		return err
	}

	_, err = tx.Exec(deleteMatchesQuery)
	if err != nil {
		log.Printf("Failed to delete all matches: %v", err)
		return err
	}
	return tx.Commit()
}

func (sqlite *SQLiteDatabase) populateMatches(rows *sql.Rows) ([]models.Match, error) {
	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var ht, at models.Team

		var homeScore, awayScore sql.NullInt64
		htAttrs := &ht.Attributes
		atAttrs := &at.Attributes

		err := rows.Scan(
			&match.ID, &match.Week, &homeScore, &awayScore, &match.IsPlayed,
			&ht.ID, &ht.Name, &htAttrs.Attack, &htAttrs.Defense, &htAttrs.Midfield, &htAttrs.HomeBoost, &ht.PlayStyle,
			&at.ID, &at.Name, &atAttrs.Attack, &atAttrs.Defense, &atAttrs.Midfield, &atAttrs.HomeBoost, &at.PlayStyle,
		)
		if err != nil {
			log.Printf("Failed to scan match row: %v", err)
			return nil, err
		}

		if homeScore.Valid {
			match.Result.HomeScore = int(homeScore.Int64)
		}
		if awayScore.Valid {
			match.Result.AwayScore = int(awayScore.Int64)
		}

		match.HomeTeam = &ht
		match.AwayTeam = &at

		matches = append(matches, match)
	}
	return matches, nil
}
