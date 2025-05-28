package database

// Reused queries are defined here

const (
	getTeamsQuery string = `
	SELECT id, name, attack, defense, midfield, home_boost, play_style FROM teams ORDER BY name;
	`

	getMatchesQuery string = `
	SELECT m.id, m.week, m.home_score, m.away_score, m.is_played,
		ht.id as home_id, ht.name as home_name, ht.attack as home_attack, ht.defense as home_defense, ht.midfield as home_midfield, ht.home_boost as home_boost, ht.play_style as home_style,
		at.id as away_id, at.name as away_name, at.attack as away_attack, at.defense as away_defense, at.midfield as away_midfield, at.home_boost as away_boost, at.play_style as away_style
	FROM matches m
	JOIN teams ht ON m.home_team_id = ht.id
	JOIN teams at ON m.away_team_id = at.id
	ORDER BY m.week, m.id;
	`

	getMatchesForWeekQuery string = `
	SELECT m.id, m.week, m.home_score, m.away_score, m.is_played,
		ht.id as home_id, ht.name as home_name, ht.attack as home_attack, ht.defense as home_defense, ht.midfield as home_midfield, ht.home_boost as home_boost, ht.play_style as home_style,
		at.id as away_id, at.name as away_name, at.attack as away_attack, at.defense as away_defense, at.midfield as away_midfield, at.home_boost as away_boost, at.play_style as away_style
	FROM matches m
	JOIN teams ht ON m.home_team_id = ht.id
	JOIN teams at ON m.away_team_id = at.id
	WHERE m.week = ?
	ORDER BY m.id;
	`

	getStateQuery string = `
	SELECT id, current_week, max_weeks FROM simulation_state WHERE id = 1;
	`

	insertMatchQuery string = `
	INSERT INTO matches (week, home_team_id, away_team_id, is_played)
	VALUES (?, ?, ?, FALSE);
	`

	updateMatchQuery string = `
	UPDATE matches
	SET home_score = ?, away_score = ?, is_played = TRUE
	WHERE id = ?;
	`

	updateWeekQuery string = `
	UPDATE simulation_state
	SET current_week = ?
	WHERE id = 1;
	`

	resetMatchesQuery string = `
	UPDATE matches
	SET home_score = NULL, away_score = NULL, is_played = FALSE;
	`

	resetStateQuery string = `
	UPDATE simulation_state
	SET current_week = 1;
	`

	deleteMatchesQuery string = `
	DELETE FROM matches;
	`
)
