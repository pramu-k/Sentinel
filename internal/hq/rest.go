package hq

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RESTServer struct {
	Store  MetricStore
	Router *gin.Engine
}

func NewRESTServer(store MetricStore) *RESTServer {
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	s := &RESTServer{
		Store:  store,
		Router: r,
	}
	s.registerRoutes()
	return s
}

func (s *RESTServer) registerRoutes() {
	s.Router.GET("/servers", s.handleListServers)
	s.Router.GET("/metrics/:server_id", s.handleGetMetrics)
	s.Router.GET("/servers/:server_id/services", s.handleGetServiceStatus)
}

func (s *RESTServer) handleListServers(c *gin.Context) {
	servers, err := s.Store.ListServers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, servers)
}

func (s *RESTServer) handleGetMetrics(c *gin.Context) {
	serverID := c.Param("server_id")
	metrics, err := s.Store.GetMetrics(c.Request.Context(), serverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func (s *RESTServer) handleGetServiceStatus(c *gin.Context) {
	serverID := c.Param("server_id")
	services, err := s.Store.GetServiceStatus(c.Request.Context(), serverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}
