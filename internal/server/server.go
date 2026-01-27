package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/handlers/auth"
	botHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/bot"

	userHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/user"
	"github.com/souravsspace/texly.chat/internal/middleware"
	authMiddleware "github.com/souravsspace/texly.chat/internal/middleware/auth"
	botRepoPkg "github.com/souravsspace/texly.chat/internal/repo/bot"

	userRepoPkg "github.com/souravsspace/texly.chat/internal/repo/user"
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
	* Repositories
	*/
	userRepo := userRepoPkg.NewUserRepo(s.db)
  botRepo := botRepoPkg.NewBotRepo(s.db)

/*
	* Handlers
	*/
	authHandler := auth.NewAuthHandler(userRepo, s.cfg)
	userHandler := userHandlerPkg.NewUserHandler(userRepo)


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
