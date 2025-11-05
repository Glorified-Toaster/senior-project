// Package controllers contains handler functions for the web application.
package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	time.Sleep(9 * time.Second) // Simulate a delay
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}
