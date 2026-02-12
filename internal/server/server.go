package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/db"
	analyticsHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/analytics"
	"github.com/souravsspace/texly.chat/internal/handlers/auth"
	botHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/bot"
	chatHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/chat"
	healthHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/health"
	publicHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/public"
	sourceHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/source"
	userHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/user"
	authMiddleware "github.com/souravsspace/texly.chat/internal/middleware/auth"
	corsMiddleware "github.com/souravsspace/texly.chat/internal/middleware/cors"
	middleware "github.com/souravsspace/texly.chat/internal/middleware/entitlement"
	rateLimitMiddleware "github.com/souravsspace/texly.chat/internal/middleware/rate_limit"
	"github.com/souravsspace/texly.chat/internal/queue"
	botRepoPkg "github.com/souravsspace/texly.chat/internal/repo/bot"
	messageRepoPkg "github.com/souravsspace/texly.chat/internal/repo/message"
	sourceRepoPkg "github.com/souravsspace/texly.chat/internal/repo/source"
	userRepoPkg "github.com/souravsspace/texly.chat/internal/repo/user"
	vectorRepoPkg "github.com/souravsspace/texly.chat/internal/repo/vector"
	"github.com/souravsspace/texly.chat/internal/services/analytics"
	billing "github.com/souravsspace/texly.chat/internal/services/billing/core"
	credits "github.com/souravsspace/texly.chat/internal/services/billing/credits"
	polar "github.com/souravsspace/texly.chat/internal/services/billing/polar"
	usage "github.com/souravsspace/texly.chat/internal/services/billing/usage"
	"github.com/souravsspace/texly.chat/internal/services/cache"
	"github.com/souravsspace/texly.chat/internal/services/chat"
	"github.com/souravsspace/texly.chat/internal/services/embedding"
	"github.com/souravsspace/texly.chat/internal/services/oauth"
	"github.com/souravsspace/texly.chat/internal/services/session"
	"github.com/souravsspace/texly.chat/internal/services/storage"
	"github.com/souravsspace/texly.chat/internal/services/vector"
	"github.com/souravsspace/texly.chat/internal/worker"
	"github.com/souravsspace/texly.chat/ui"
	"github.com/souravsspace/texly.chat/widget"
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
	* Redis Client & Cache Service
	 */
	redisClient := db.GetRedisClient()
	cacheService := cache.NewCacheService(redisClient)

	/*
	* Billing Services
	 */
	// Core usage tracking
	usageService := usage.NewUsageService(s.db)
	
	// Supporting services for billing cycle management
	creditsService := credits.NewCreditsService(s.db)
	polarService := polar.NewPolarService(s.cfg)
	billingCycleService := billing.NewBillingCycleService(s.db, creditsService, polarService)

	/*
	* Repositories
	 */
	userRepo := userRepoPkg.NewUserRepo(s.db, cacheService)
	botRepo := botRepoPkg.NewBotRepo(s.db, cacheService)
	sourceRepo := sourceRepoPkg.NewSourceRepo(s.db, cacheService)
	messageRepo := messageRepoPkg.New(s.db)

	/*
	* Queue and Worker
	 */
	jobQueue := queue.NewInMemoryQueue(100, 3) // buffer: 100, workers: 3

	/*
	* OAuth Services
	 */
	oauthStateService := oauth.NewStateService(redisClient)
	oauthService := oauth.NewOAuthService(s.cfg, s.db)

	// Initialize MinIO storage service
	storageService, err := storage.NewMinIOStorageService(
		s.cfg.MinIOEndpoint,
		s.cfg.MinIOAccessKey,
		s.cfg.MinIOSecretKey,
		s.cfg.MinIOBucket,
		s.cfg.MinIOUseSSL,
		s.cfg.MaxUploadSizeMB,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO storage service: %w", err)
	}
	fmt.Println("✅ MinIO storage service initialized")

	// Initialize embedding service, vector search, and chat service if API key is configured
	var embeddingService *embedding.EmbeddingService
	var vectorRepo *vectorRepoPkg.VectorRepository
	var searchService *vector.SearchService
	var chatService *chat.ChatService

	if s.cfg.OpenAIAPIKey != "" {
		embeddingService = embedding.NewEmbeddingService(
			s.cfg.OpenAIAPIKey,
			s.cfg.EmbeddingModel,
			s.cfg.EmbeddingDimension,
		)
		vectorRepo = vectorRepoPkg.NewVectorRepository(s.db)
		searchService = vector.NewSearchService(s.db, vectorRepo, embeddingService)
		chatService = chat.NewChatService(
			embeddingService,
			searchService,
			messageRepo,
			s.cfg.ChatModel,
			s.cfg.ChatTemperature,
			s.cfg.MaxContextChunks,
			s.cfg.OpenAIAPIKey,
		)
		fmt.Println("✅ Embedding service initialized")
		fmt.Println("✅ Vector search service initialized")
		fmt.Println("✅ Chat service initialized")
	} else {
		fmt.Println("⚠️  OpenAI API key not configured - vector embeddings and chat disabled")
	}

	workerInstance := worker.NewWorker(s.db, embeddingService, vectorRepo, storageService, sourceRepo, botRepo, usageService)

	// Start worker pool
	jobQueue.Start(ctx, workerInstance.ProcessScrapeJob)

	// Start daily billing worker
	go worker.StartDailyBillingJob(ctx, billingCycleService)
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
	googleHandler := auth.NewGoogleHandler(oauthService, oauthStateService, s.cfg)
	userHandler := userHandlerPkg.NewUserHandler(userRepo)
	sourceHandler := sourceHandlerPkg.NewSourceHandler(sourceRepo, botRepo, jobQueue, storageService, usageService, s.cfg.MaxUploadSizeMB)
	analyticsService := analytics.NewAnalyticsService(messageRepo)
	analyticsHandler := analyticsHandlerPkg.NewAnalyticsHandler(analyticsService)

	/*
	* Health Endpoints
	 */
	healthHandler := healthHandlerPkg.NewHealthHandler(s.db, redisClient)

	/*
	* Middleware
	 */
	if os.Getenv("ENVIRONMENT") == "development" {
		s.engine.Use(corsMiddleware.CORS())
	}

	// Initialize rate limiter
	rateLimiter, err := rateLimitMiddleware.NewRateLimiter(s.cfg, s.db)
	if err != nil {
		fmt.Printf("⚠️  Failed to initialize rate limiter: %v\n", err)
		fmt.Println("⚠️  Continuing without rate limiting")
	} else {
		fmt.Println("✅ Rate limiter initialized")
	}

	// Initialize entitlement middleware
	entitlementMiddleware := middleware.NewEntitlementMiddleware(s.db)

	/*
	* Routes
	 */
	// Health routes (no auth required)
	s.engine.GET("/health", healthHandler.GetBasicHealth)
	s.engine.GET("/health/db", healthHandler.GetDBHealth)
	s.engine.GET("/health/redis", healthHandler.GetRedisHealth)

	apiGroup := s.engine.Group("/api")
	{
		/*
		* Auth routes
		 */
		apiGroup.POST("/auth/signup", authHandler.Signup)

		apiGroup.POST("/auth/login", authHandler.Login)
		apiGroup.GET("/auth/google", googleHandler.GoogleLogin)
		apiGroup.GET("/auth/google/callback", googleHandler.GoogleCallback)

		/*
		* User routes
		 */
		apiGroup.GET("/users/me", authMiddleware.Auth(s.cfg), userHandler.GetMe)

		/*
		* Bot routes
		 */
		botHandler := botHandlerPkg.NewBotHandler(botRepo)
		apiGroup.POST("/bots", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitBotCreation), botHandler.CreateBot)
		apiGroup.GET("/bots", authMiddleware.Auth(s.cfg), botHandler.ListBots)
		apiGroup.GET("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.GetBot)
		apiGroup.PUT("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.UpdateBot)
		apiGroup.DELETE("/bots/:id", authMiddleware.Auth(s.cfg), botHandler.DeleteBot)

		/*
		* Source routes (nested under bots)
		 */
		apiGroup.POST("/bots/:id/sources", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitSourceCreation), sourceHandler.CreateSource)                // URL source
		apiGroup.POST("/bots/:id/sources/upload", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitStorage), sourceHandler.UploadFileSource)        // File upload - enforcing storage limit (placeholder)
		apiGroup.POST("/bots/:id/sources/text", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitSourceCreation), sourceHandler.CreateTextSource)       // Text source
		apiGroup.POST("/bots/:id/sources/sitemap", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitSourceCreation), sourceHandler.CreateSitemapSource) // Sitemap crawl
		apiGroup.GET("/bots/:id/sources", authMiddleware.Auth(s.cfg), sourceHandler.ListSources)
		apiGroup.GET("/bots/:id/sources/:sourceId", authMiddleware.Auth(s.cfg), sourceHandler.GetSource)
		apiGroup.DELETE("/bots/:id/sources/:sourceId", authMiddleware.Auth(s.cfg), sourceHandler.DeleteSource)

		/*
		* Chat routes
		 */
		chatHandler := chatHandlerPkg.NewChatHandler(botRepo, chatService, usageService)
		apiGroup.POST("/bots/:id/chat", authMiddleware.Auth(s.cfg), entitlementMiddleware.EnforceLimit(middleware.LimitMessageSend), chatHandler.StreamChat)

		/*
		* Analytics routes
		 */
		apiGroup.GET("/analytics/bots/:id", authMiddleware.Auth(s.cfg), analyticsHandler.GetBotAnalytics)
		apiGroup.GET("/analytics/bots/:id/daily", authMiddleware.Auth(s.cfg), analyticsHandler.GetBotDailyStats)
		apiGroup.GET("/analytics/user", authMiddleware.Auth(s.cfg), analyticsHandler.GetUserAnalytics)
		apiGroup.GET("/analytics/sessions/:id/messages", authMiddleware.Auth(s.cfg), analyticsHandler.GetSessionMessages)
	}

	/*
	* Public API routes for widget
	 */
	sessionService := session.NewSessionService()
	publicHandler := publicHandlerPkg.NewPublicHandler(botRepo, sessionService, chatService)

	publicGroup := s.engine.Group("/api/public")
	publicGroup.Use(corsMiddleware.WidgetCORS(botRepo))
	// Apply rate limiting to public endpoints
	if rateLimiter != nil {
		publicGroup.Use(rateLimiter.PublicRateLimitMiddleware())
	}
	{
		// Widget configuration
		publicGroup.GET("/bots/:id/config", publicHandler.GetWidgetConfig)

		// Session management
		publicGroup.POST("/chats", publicHandler.CreateSession)

		// Chat streaming
		publicGroup.POST("/chats/:session_id/messages", publicHandler.StreamChatPublic)
	}

	/*
	* Serve Frontend
	 */
	if err := ui.RegisterRoutes(s.engine); err != nil {
		fmt.Printf("Warning: Failed to register web routes: %v\n", err)
	}

	/*
	* Serve Widget
	 */
	if err := widget.RegisterRoutes(s.engine); err != nil {
		fmt.Printf("Warning: Failed to register widget routes: %v\n", err)
	}

	addr := fmt.Sprintf(":%s", s.cfg.Port)
	return s.engine.Run(addr)
}
