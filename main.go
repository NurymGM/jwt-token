package main

import (
	"net/http"

	"github.com/NurymGM/jwt-token/controllers"
	"github.com/NurymGM/jwt-token/initializers"
	"github.com/NurymGM/jwt-token/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello from root route!",
		})
	})

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.LogIn)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
