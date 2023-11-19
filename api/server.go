package api

import (
	db "github.com/LKarrie/mdc-server/db/sqlc"
	"github.com/LKarrie/mdc-server/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server servers HTTP requests for our banking service.
type Server struct {
	config util.Config
	docker *Docker
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, d *Docker, store db.Store) (*Server, error) {

	server := &Server{
		config: config,
		docker: d,
		store:  store,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	
	router := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true
	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS","GET","POST")
	corsConfig.ExposeHeaders = []string{"Content-Disposition"}
	// Register the middleware
	router.Use(cors.New(corsConfig))

	router.GET("/images/list", server.listImage)
	router.POST("/images/pull", server.pullImage)
	router.POST("/images/pull/auth", server.pullImageWithAuth)
	router.POST("/images/create/tag", server.tagImage)
	router.POST("/images/save", server.saveImage)
	router.POST("/images/load", server.loadImage)
	router.POST("/images/push", server.pushImage)
	router.POST("/images/push/auth", server.pushImageWithAuth)
	// add routers to router
	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
