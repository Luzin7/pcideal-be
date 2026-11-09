package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Luzin7/pcideal-be/infra/database"
	"github.com/Luzin7/pcideal-be/infra/external"
	"github.com/Luzin7/pcideal-be/infra/http/controllers"
	partControllers "github.com/Luzin7/pcideal-be/infra/http/controllers/part"
	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/Luzin7/pcideal-be/infra/http/routes"
	"github.com/Luzin7/pcideal-be/infra/repositories"
	"github.com/Luzin7/pcideal-be/internal/useCases/buildAttempt"
	"github.com/Luzin7/pcideal-be/internal/useCases/part"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	indexTimeout    = 30 * time.Second
	shutdownTimeout = 10 * time.Second
)

func main() {
	isProduction := os.Getenv("ENV") == "production"

	if !isProduction {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	if isProduction && os.Getenv("TRUSTED_PROXIES") == "" && os.Getenv("CLIENT_IP_HEADER") == "" {
		log.Println("WARNING: neither TRUSTED_PROXIES nor CLIENT_IP_HEADER is set in production; every client behind the proxy shares one rate-limit bucket. Set CLIENT_IP_HEADER=CF-Connecting-IP behind Cloudflare, or TRUSTED_PROXIES to the proxy CIDRs.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	clientIPResolver, clientIPErr := middlewares.NewClientIPResolver(os.Getenv("CLIENT_IP_HEADER"))
	if clientIPErr != nil {
		log.Fatalf("Erro ao configurar CLIENT_IP_HEADER: %v", clientIPErr)
	}

	connectionString := os.Getenv("DATABASE_URL")
	databaseName := os.Getenv("PCIDEAL_DB_NAME")

	db, err := database.MongoConnection(connectionString, databaseName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %s", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := db.Client().Disconnect(ctx); err != nil {
			log.Printf("Erro ao desconectar do MongoDB: %v", err)
		}
	}()

	indexCtx, cancelIndexes := context.WithTimeout(context.Background(), indexTimeout)
	err = database.EnsureIndexes(indexCtx, db)
	cancelIndexes()
	if err != nil {
		log.Fatalf("Erro ao garantir índices do MongoDB: %v", err)
	}

	indexCheckCtx, cancelIndexCheck := context.WithTimeout(context.Background(), indexTimeout)
	uniqueIndexReady, indexCheckErr := database.UniqueStoreURLIndexReady(indexCheckCtx, db)
	cancelIndexCheck()
	if indexCheckErr != nil {
		log.Printf("[database] ERROR: could not verify unique index %s: %v", database.StoreURLIndexName, indexCheckErr)
	} else if !uniqueIndexReady {
		log.Printf(
			"[database] ERROR: unique index %s is missing; imports are NOT idempotent. Run the duplicate audit: %s",
			database.StoreURLIndexName, database.StoreURLDuplicateAuditCommand,
		)
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
	getBuildRecsController := partControllers.NewGetBuildRecommendationsController(generateBuildRecsUC, buildAttemptService, clientIPResolver)
	healthController := controllers.NewHealthController(db.Client(), uniqueIndexReady)

	router, err := routes.SetupRouter(routes.Dependencies{
		GetAllPartsController:  getAllPartsController,
		GetPartByIDController:  getPartByIDController,
		GetBuildRecsController: getBuildRecsController,
		HealthController:       healthController,
		ClientIPResolver:       clientIPResolver,
	})
	if err != nil {
		log.Fatalf("Erro ao configurar rotas: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	log.Printf("Servidor iniciando na porta %s...", port)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Erro ao subir o servidor: %v", err)
		}
	case <-signalCtx.Done():
		log.Println("Encerrando servidor...")

		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Erro ao encerrar o servidor: %v", err)
		}
	}
}
