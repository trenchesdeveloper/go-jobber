package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/trenchesdeveloper/jobber/config"
	db "github.com/trenchesdeveloper/jobber/internal/db/sqlc"
	"go.uber.org/zap"
)

type server struct {
	config *config.AppConfig
	logger *zap.SugaredLogger
	store  db.Store
}

func (s *server) mount() http.Handler {
	r := gin.Default()

	// ping route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Auth routes
	auth := r.Group("/auth")
	auth.POST("/register", s.Register)

	return r
}

func (s *server) start(mux http.Handler) error {
	srv := &http.Server{
		Addr:         s.config.ServerPort,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	s.logger.Infow("Starting server", "port", s.config.ServerPort, "env", s.config.Environment)
	return srv.ListenAndServe()
}
