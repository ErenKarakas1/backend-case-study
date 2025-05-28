# Insider Backend Case Study

**Main prompt:** In this project, we expect you to complete a simulation. In this simulation, there will be a group of
football teams and the simulation will show match results and the league table. Your task is to
estimate the final league table.

## Table of Contents

- [Project Structure](#project-structure)
- [Endpoints](#endpoints)
- [Environment Variables](#environment-variables)
- [Usage](#usage)

## Project Structure

```
/root
    database/
        database.go             # Collection of database operations
        league.db
        queries.go              # Collection of reused queries
    handlers/
        handlers_test.go
        handlers.go
    http_templates              # Collection of example HTTP request templates
    models/                     # Project-wide used types are defined here
        league.go
        match.go
        team.go
    services/
        leaguePredictor_test.go
        leaguePredictor.go
        leagueService.go
        leagueTable_test.go
        leagueTable.go
        matchScheduler_test.go
        matchScheduler.go
        matchSimulator_test.go
        matchSimulator.go
        services.go
    templates/
        index.html
    go.mod
    go.sum
    main.go
    README.md
```

## Endpoints

- **GET /api/simulation**

Return the full current state of the simulation.

```json
{
    "current_week": int,
    "max_weeks": int,
    "table": [
        {
            "position": int,
            "team": {
                "id": int,
                "name": "string"
            },
            "played": int,
            "won": int,
            "drawn": int,
            "lost": int,
            "goals_for": int,
            "goals_against": int,
            "goal_diff": int,
            "points": int
        }
        // ...
    ],
    "matches": [
        {
            "id": int,
            "week": int,
            "home_team": {
                "id": int,
                "name": "string"
            },
            "away_team": {
                "id": int,
                "name": "string"
            },
            "result": {
                "home_score": int,
                "away_score": int
            },
            "is_played": boolean
        },
        // ...
    ],
    // if after 4th week
    "championship_odds": [
        {
            "team_id": int,
            "team_name": "string",
            "probability": float
        },
        // 3 more
    ]
    // endif
}
```

- **POST /api/simulation/next-week**

Simulate next week's fixtures. Returns only that week's results.

```json
{
    "played_week": int,
    "matches": [
        {
            "id": int,
            "week": int,
            "home_team": {
                "id": int,
                "name": "string"
            },
            "away_team": {
                "id": int,
                "name": "string"
            },
            "result": {
                "home_score": int,
                "away_score": int
            },
            "is_played": boolean
        },
        // 1 more
    ]
}
```

- **POST /api/simulation/remaining-weeks**

Simulate all remaining weeks in one go. Returns the full simulation state in the same format as **GET /api/simulation**.

- **POST /api/simulation/reset**

Reset the simulation back to week 1 and reshuffle the schedule. On success, returns:

```json
{
    "message": "Simulation reset successfully"
}
```

- **PUT /api/simulation/edit-match-result**

Edit a past match's score. Payload:

```http
Content-Type: application/json

{
  "match_id": int,
  "home_score": int,
  "away_score": int
}
```

Returns the resulting state of the simulation, both league table and championship odds are recalculated. Its in the same format as **GET /api/simulation**.

## Environment Variables (.env.example)

```env

# Load GIN in debug mode
GIN_MODE=debug # or release for production

# Database file path
DATABASE_URL=database/league.db

```

## Usage

For local development:
1. Run `go run .` in the root directory
2. The server will run on `http://localhost:8080` by default.

Alternatively:


`Simulate Next Week` button will send a request to the `api/simulation/next-week` endpoint.

`Simulate Remaining Weeks` button will send a request to the `api/simulation/remaining-weeks` endpoint.

`Reset Simulation` button will send a request to the `api/simulation/reset` endpoint.

After Week 4, the frontend will also display the predicted odds for each team to win the championship. The list is ordered with respect to current positions of the teams.

You can send a request to the `api/simulation/edit-match-result` endpoint to edit a match's result. The frontend will reflect the new standings, history and odds after refreshing the page.
