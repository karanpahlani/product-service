package httpserver

import (
	"log"

	"github.com/gin-gonic/gin"

	"product-service/internal/database"
	"product-service/internal/handlers"
	"product-service/internal/repository"
	"product-service/internal/service"
)

type Server struct {
	router  *gin.Engine
	handler *handlers.ProductHandler
}

func NewServer() (*Server, error) {
	db, err := database.NewDynamoDBClient()
	if err != nil {
		return nil, err
	}

	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	handler := handlers.NewProductHandler(svc)

	router := gin.Default()

	server := &Server{
		router:  router,
		handler: handler,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	
	api.GET("/health", s.handler.HealthCheck)
	
	products := api.Group("/products")
	{
		products.POST("", s.handler.CreateProduct)
		products.GET("", s.handler.GetAllProducts)
		products.GET("/category", s.handler.GetProductsByCategory)
		products.GET("/:id", s.handler.GetProduct)
		products.PUT("/:id", s.handler.UpdateProduct)
		products.DELETE("/:id", s.handler.DeleteProduct)
	}
}

func (s *Server) Run(addr string) error {
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}

