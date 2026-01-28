package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/handlers/auth"
	botHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/bot"
	sourceHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/source"
	userHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/user"
	"github.com/souravsspace/texly.chat/internal/middleware"
	authMiddleware "github.com/souravsspace/texly.chat/internal/middleware/auth"
	"github.com/souravsspace/texly.chat/internal/queue"
	botRepoPkg "github.com/souravsspace/texly.chat/internal/repo/bot"
	sourceRepoPkg "github.com/souravsspace/texly.chat/internal/repo/source"
	userRepoPkg "github.com/souravsspace/texly.chat/internal/repo/user"
	vectorRepoPkg "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/worker"
	"github.com/souravsspace/texly.chat/ui"
	"gorm.io/gorm"
)

/*
* Server holds the dependencies for the HTTP server
 */
type Server struct {
	engine *gin.Engine
	db     *gorm.DB
	cfg    configs.Config
}

/*
* New creates a new Server instance
*/
func New(db *gorm.DB, cfg configs.Config) *Server {
	return &Server{
		engine: gin.Default(),
		db:     db,
		cfg:    cfg,
	}
}

/*
* Run configures routes and starts the HTTP server
*/
func (s *Server) Run() error {
	/*
	* Setup context for graceful shutdown
	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/*
	* Repositories
	*/
	userRepo := userRepoPkg.NewUserRepo(s.db)
	botRepo := botRepoPkg.NewBotRepo(s.db)
	sourceRepo := sourceRepoPkg.NewSourceRepo(s.db)

	/*
	* Queue and Worker
	*/
	jobQueue := queue.NewInMemoryQueue(100, 3) // buffer: 100, workers: 3
	
	// Initialize embedding service and vector repository if API key is configured
	var embeddingService *embedding.EmbeddingService
	var vectorRepo *vectorRepoPkg.VectorRepository
	
	if s.cfg.OpenAIAPIKey != "" {
		embeddingService = embedding.NewEmbeddingService(
			s.cfg.OpenAIAPIKey,
			s.cfg.EmbeddingModel,
			s.cfg.EmbeddingDimension,
		)
		vectorRepo = vectorRepoPkg.NewVectorRepository(s.db)
		fmt.Println("✅ Embedding service initialized")
	} else {
		fmt.Println("⚠️  OpenAI API key not configured - vector embeddings disabled")
	}
	
	workerInstance := worker.NewWorker(s.db, embeddingService, vectorRepo)
	
	// Start worker pool
	jobQueue.Start(ctx, workerInstance.ProcessScrapeJob)
	fmt.Println("✅ Worker pool started")

	// Setup graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\n🛑 Shutdown signal received, stopping workers...")
		cancel()
		jobQueue.Stop()
		os.Exit(0)
	}()

	/*
	* Handlers
	*/
	authHandler := auth.NewAuthHandler(userRepo, s.cfg)
	userHandler := userHandlerPkg.NewUserHandler(userRepo)
	sourceHandler := sourceHandlerPkg.NewSourceHandler(sourceRepo, botRepo, jobQueue)


/*
	* Middleware
	*/
	if os.Getenv("ENVIRONMENT") == "development" {
		s.engine.Use(middleware.CORS())
	}

/*
	* Routes
	*/
	apiGroup := s.engine.Group("/api")
	{
/*
		* Auth routes
		*/
		apiGroup.POST("/auth/signup", authHandler.Signup)
		apiGroup.POST("/auth/login", authHandler.Login)

/*
		* User routes
		*/
		apiGroup.GET("/users/me", authMiddleware.Auth(s.cfg), userHandler.GetMe)


		/*
		* Bot routes
		*/
		botHandler := botHandlerPkg.NewBotHandler(botRepo)
		apiGroup.POST("/bots", authMiddleware.Auth(s.cfg), botHandler.CreateBot)
		apiGroup.GET("/bots", authMiddleware.Auth(s.cfg), botHandler.ListBots)
		apiGroup.GET("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.GetBot)
		apiGroup.PUT("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.UpdateBot)
		apiGroup.DELETE("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.DeleteBot)

		/*
		* Source routes (nested under bots)
		*/
		apiGroup.POST("/bots/:id/sources", authMiddleware.Auth(s.cfg), sourceHandler.CreateSource)
		apiGroup.GET("/bots/:id/sources", authMiddleware.Auth(s.cfg), sourceHandler.ListSources)
		apiGroup.GET("/bots/:id/sources/:sourceId", authMiddleware.Auth(s.cfg), sourceHandler.GetSource)
		apiGroup.DELETE("/bots/:id/sources/:sourceId", authMiddleware.Auth(s.cfg), sourceHandler.DeleteSource)
	}

/*
	* Serve Frontend
	*/
	if err := ui.RegisterRoutes(s.engine); err != nil {
		fmt.Printf("Warning: Failed to register web routes: %v\n", err)
	}

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.engine.Run(addr)
}
