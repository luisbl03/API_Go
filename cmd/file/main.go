package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

func main() {
	config.Load("config/config.toml")
	api := gin.Default()
	api.POST("/:username/:doc_id", upload)
	api.GET("/:username/:doc_id", get)
	api.Run("127.0.0.1:8082")
	
}

func upload(c *gin.Context) {
	var json models.Json
	err := c.BindJSON(&json)
	if err != nil {
		c.JSON(400, gin.H{"error":"invalid json"})
		return
	}
	username := c.Param("username")
	doc_id := c.Param("doc_id")
	status := api.Upload(username, json, doc_id)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(201, gin.H{"size":status})

}

func get(c *gin.Context) {
	username := c.Param("username")
	doc_id := c.Param("doc_id")
	json, status := api.GetFile(doc_id, username)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(200, json)
}

func Status(status int) (string, int) {
	if status == constants.ERROR {
		return "internal error", 500
	}
	if status == constants.EXISTS {
		return "exists", 409
	}
	if status == constants.NOT_FOUND {
		return "not found", 404
	}
	if status == constants.CREATED {
		return "",201
	}
	if status == constants.INVALID_PASSWORD {
		return "invalid password", 401
	}
	return "", 200
}