package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) Register(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"user": "registered"})
}
