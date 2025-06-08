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
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("PCIDEAL_DB_NAME")
	db, err := database.MongoConnection(connectionString, databaseName)

	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados")
	}

	googleAiClient, err := external.NewGoogleAIClient(os.Getenv("GOOGLE_AI_API_KEY"), db)
	if err != nil {
		log.Fatal("Erro ao criar o cliente Google AI:", err)
	}

	partRepository := repositories.NewPartRepository(db)
	scraperClient := external.NewScraperHTTPClient(os.Getenv("SCRAPER_API_URL"))
	partMatchingService := matching.NewPartMatchingService(db)
	partService := services.NewPartService(partRepository, scraperClient, googleAiClient, partMatchingService)
	partController := controllers.NewPartController(partService)

	router := routes.SetupRouter(partController)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao subir o servidor:", err)
	}
}
