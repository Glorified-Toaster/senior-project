// Package controllers contains handler functions for the web application.
package controllers

import (
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}
