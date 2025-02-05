package main

import (
	"net/http"

	"github.com/NurymGM/jwt-token/initializers"
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
			"message": "pong3",
		})
	})

	r.Run()
}
