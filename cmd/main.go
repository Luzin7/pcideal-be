package main

import (
	"log"
	"os"

	"github.com/Luzin7/pcideal-be/infra/database"
	"github.com/Luzin7/pcideal-be/infra/external"
	"github.com/Luzin7/pcideal-be/infra/repositories"
	"github.com/Luzin7/pcideal-be/internal/domain/matching"
	"github.com/Luzin7/pcideal-be/internal/http/controllers"
	"github.com/Luzin7/pcideal-be/internal/http/routes"
	"github.com/Luzin7/pcideal-be/internal/http/services"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	connectionString := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("PCIDEAL_DB_NAME")
	db, err := database.MongoConnection(connectionString, databaseName)

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %s", err)
	}

	googleAiClient, err := external.NewGoogleAIClient(os.Getenv("GOOGLE_AI_API_KEY"), db)
	if err != nil {
		log.Fatal("Erro ao criar o cliente Google AI:", err)
	}

	partRepository := repositories.NewPartRepository(db)
	buildAttemptRepository := repositories.NewBuildAttemptRepository(db)
	scraperClient := external.NewScraperHTTPClient(os.Getenv("SCRAPER_API_URL"), os.Getenv("SCRAPER_API_KEY"))
	partMatchingService := matching.NewPartMatchingService(db)
	buildAttemptService := services.NewBuildAttemptService(buildAttemptRepository)
	partService := services.NewPartService(partRepository, scraperClient, googleAiClient, partMatchingService)
	partController := controllers.NewPartController(partService, buildAttemptService)

	router := routes.SetupRouter(partController)

	log.Printf("Servidor iniciando na porta %s...", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Erro ao subir o servidor:", err)
    }
}
