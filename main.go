package main

import (
	"log"
	"os"

	"insider/database"
	"insider/handlers"
	"insider/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "database/league.db"
	}

	db := database.NewSQLiteDatabase(dbPath)
	db.Initialize()
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Failed to close database: ", err)
		}
	}()

	log.SetFlags(log.Llongfile)

	teams, err := db.GetTeams()
	if err != nil {
		log.Fatal("Failed to retrieve teams from database: ", err)
	}

	simulator := services.NewMatchSimulator()
	table := services.NewLeagueTable(teams)
	scheduler := services.NewMatchScheduler()
	predictor := services.NewLeaguePredictor(simulator, table)
	leagueService := services.NewLeagueService(db, simulator, table, scheduler, predictor)

	if err := leagueService.ResetSimulation(); err != nil {
		log.Fatal("Failed to initialize league simulation: ", err)
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.LoadHTMLGlob("templates/*")

	router.GET("/api/simulation", handlers.GetSimulationState(leagueService))
	router.POST("/api/simulation/next-week", handlers.SimulateNextWeek(leagueService))
	router.POST("/api/simulation/remaining-weeks", handlers.SimulateRemainingWeeks(leagueService))
	router.POST("/api/simulation/reset", handlers.ResetSimulation(leagueService))
	router.PUT("/api/simulation/edit-match-result", handlers.EditMatchResult(leagueService))

	router.GET("/", handlers.ServeIndex())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server is running on http://localhost:" + port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
