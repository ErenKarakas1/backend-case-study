package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"insider/database"
	"insider/handlers"
	"insider/models"
	"insider/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	db := database.NewSQLiteDatabase(":memory:")
	db.Initialize()

	teams, err := db.GetTeams()
	if err != nil {
		t.Fatalf("Failed to retrieve teams from database: %v", err)
	}

	scheduler := services.NewMatchScheduler()
	simulator := services.NewMatchSimulator()
	table := services.NewLeagueTable(teams)
	predictor := services.NewLeaguePredictor(simulator, table)
	svc := services.NewLeagueService(db, simulator, table, scheduler, predictor)

	if err := svc.ResetSimulation(); err != nil {
		t.Fatalf("ResetSimulation: %v", err)
	}

	r := gin.New()
	sim := r.Group("/api/simulation")
	{
		sim.GET("", handlers.GetSimulationState(svc))
		sim.POST("/next-week", handlers.SimulateNextWeek(svc))
		sim.POST("/remaining-weeks", handlers.SimulateRemainingWeeks(svc))
		sim.POST("/reset", handlers.ResetSimulation(svc))
	}
	return r
}

func TestIntegration_GetSimulationState(t *testing.T) {
	router := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/simulation", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.LeagueSimulation
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 1, resp.CurrentWeek)
	assert.Equal(t, 6, resp.MaxWeeks)
	assert.Len(t, resp.Table, 4)
	assert.Len(t, resp.Matches, 12)
	assert.Empty(t, resp.ChampionshipOdds)
}

func TestIntegration_SimulateNextWeek(t *testing.T) {
	router := setupTestRouter(t)

	// advance one week
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/simulation/next-week", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.WeekSimulation
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.PlayedWeek)

	played := 0
	for _, m := range resp.Matches {
		if m.Week == 1 && m.IsPlayed {
			played++
		}
	}
	assert.Equal(t, 2, played) // two matches week 1
}

func TestIntegration_SimulateFourthWeek(t *testing.T) {
	router := setupTestRouter(t)

	// advance to week 4
	for range 3 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/simulation/next-week", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	// now simulate week 4
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/simulation/next-week", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.WeekSimulation
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, 4, resp.PlayedWeek)

	played := 0
	for _, m := range resp.Matches {
		if m.Week == 4 && m.IsPlayed {
			played++
		}
	}

	assert.Equal(t, 2, played) // two matches week 4

	// check championship odds are now available
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/simulation", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var leagueResp models.LeagueSimulation
	json.Unmarshal(w.Body.Bytes(), &leagueResp)

	assert.Equal(t, 5, leagueResp.CurrentWeek)

	assert.NotEmpty(t, leagueResp.ChampionshipOdds, "should have championship odds after week 4")
}

func TestIntegration_SimulateRemainingWeeks(t *testing.T) {
	router := setupTestRouter(t)

	// simulate all weeks
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/simulation/remaining-weeks", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp models.LeagueSimulation
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.True(t, resp.CurrentWeek > resp.MaxWeeks)

	for _, m := range resp.Matches {
		assert.True(t, m.IsPlayed)
	}
}

func TestIntegration_ResetSimulation(t *testing.T) {
	router := setupTestRouter(t)

	httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/api/simulation/next-week", nil)
	router.ServeHTTP(httptest.NewRecorder(), req1)

	w := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/simulation/reset", nil)
	router.ServeHTTP(w, req2)

	assert.Equal(t, 200, w.Code)
}
