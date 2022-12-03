package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ip-05/quizzus/controllers"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	auth := new(controllers.AuthController)

	authGroup := router.Group("auth")

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	authGroup.POST("/login", auth.Login)

	authGroup.GET("/login", auth.Test)

	return router

}
