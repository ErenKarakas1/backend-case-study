package handlers

import (
	"net/http"

	"insider/services"

	"github.com/gin-gonic/gin"
)

// GetSimulationState returns the current state of the league simulation
func GetSimulationState(service services.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := service.GetCurrentState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
	}
}

// SimulateNextWeek simulates the next week of matches
func SimulateNextWeek(service services.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := service.SimulateNextWeek()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
	}
}

// SimulateRemainingWeeks simulates all remaining weeks of the league
func SimulateRemainingWeeks(service services.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := service.SimulateRemainingWeeks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, state)
	}
}

// ResetSimulation resets the entire simulation
func ResetSimulation(service services.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := service.ResetSimulation()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Simulation reset successfully"})
	}
}

// EditMatchResult allows editing the result of a match
func EditMatchResult(service services.LeagueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			MatchID   int `json:"match_id" binding:"required"`
			HomeScore int `json:"home_score" binding:"required"`
			AwayScore int `json:"away_score" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		err := service.UpdateMatchResult(req.MatchID, req.HomeScore, req.AwayScore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sim, err := service.GetCurrentState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sim)
	}
}

// ServeIndex serves the main HTML page
func ServeIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}
