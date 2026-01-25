package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/souravsspace/texly.chat/configs"
	"github.com/souravsspace/texly.chat/internal/handlers/auth"
	postHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/post"
	userHandlerPkg "github.com/souravsspace/texly.chat/internal/handlers/user"
	"github.com/souravsspace/texly.chat/internal/middleware"
	authMiddleware "github.com/souravsspace/texly.chat/internal/middleware/auth"
	postRepoPkg "github.com/souravsspace/texly.chat/internal/repo/post"
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
	postRepo := postRepoPkg.NewPostRepo(s.db)

/*
	* Handlers
	*/
	authHandler := auth.NewAuthHandler(userRepo, s.cfg)
	userHandler := userHandlerPkg.NewUserHandler(userRepo)
	postHandler := postHandlerPkg.NewPostHandler(postRepo)

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
		* Post routes
		*/
		apiGroup.POST("/posts", authMiddleware.Auth(s.cfg), postHandler.CreatePost)
		apiGroup.GET("/posts", postHandler.GetPosts)
		apiGroup.GET("/posts/:id", postHandler.GetPost)
		apiGroup.PUT("/posts/:id", authMiddleware.Auth(s.cfg), postHandler.UpdatePost)
		apiGroup.DELETE("/posts/:id", authMiddleware.Auth(s.cfg), postHandler.DeletePost)
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
