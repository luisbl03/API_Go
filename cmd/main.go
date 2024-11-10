package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	api := gin.Default()
	api.GET("/version", getVersion) 
	api.Run(":8080")

}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	c.JSON(200, gin.H{"version":"1.0.0"})
}