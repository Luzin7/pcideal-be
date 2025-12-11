package main

import (
	"log"
	"os"

	"github.com/Luzin7/pcideal-be/infra/database"
	"github.com/Luzin7/pcideal-be/infra/external"
	partControllers "github.com/Luzin7/pcideal-be/infra/http/controllers/part"
	"github.com/Luzin7/pcideal-be/infra/http/routes"
	"github.com/Luzin7/pcideal-be/infra/repositories"
	"github.com/Luzin7/pcideal-be/internal/useCases/buildAttempt"
	"github.com/Luzin7/pcideal-be/internal/useCases/part"
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
	buildAttemptService := buildAttempt.NewBuildAttemptService(buildAttemptRepository)

	updatePartUC := part.NewUpdatePartUseCase(partRepository, scraperClient)
	updatePartsUC := part.NewUpdatePartsUseCase(partRepository, scraperClient)
	getAllPartsUC := part.NewGetAllPartsUseCase(partRepository)
	getPartByIDUC := part.NewGetPartByIDUseCase(partRepository, updatePartUC)

	selectBestCPUUC := part.NewSelectBestCPUUseCase(updatePartsUC)
	selectBestGPUUC := part.NewSelectBestGPUUseCase(updatePartsUC)
	selectBestPSUUC := part.NewSelectBestPSUUseCase(updatePartsUC)
	selectBestRAMUC := part.NewSelectBestRAMUseCase(updatePartsUC)
	selectBestMOBOUC := part.NewSelectBestMOBOUseCase(updatePartsUC)
	selectBestSSDUC := part.NewSelectBestSSDUseCase(updatePartsUC)

	generateBuildRecsUC := part.NewGenerateBuildRecommendationsUseCase(
		partRepository,
		scraperClient,
		googleAiClient,
		updatePartsUC,
		selectBestCPUUC,
		selectBestGPUUC,
		selectBestPSUUC,
		selectBestRAMUC,
		selectBestMOBOUC,
		selectBestSSDUC,
	)

	getAllPartsController := partControllers.NewGetAllPartsController(getAllPartsUC)
	getPartByIDController := partControllers.NewGetPartByIDController(getPartByIDUC)
	getBuildRecsController := partControllers.NewGetBuildRecommendationsController(generateBuildRecsUC, buildAttemptService)

	router := routes.SetupRouter(getAllPartsController, getPartByIDController, getBuildRecsController)

	log.Printf("Servidor iniciando na porta %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Erro ao subir o servidor:", err)
	}
}
