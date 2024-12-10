package main

import (
	"log"
	"github.com/luideoz/API_Go/constants"
	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/models"
)

func main() {
	api := gin.Default()
	api.POST("/:username/:doc_id", upload)
	api.Run(":8083")
}

func upload(c *gin.Context) {
	username := c.Param("username")
	doc_id := c.Param("doc_id")
	var json models.Json
	err := c.BindJSON(&json)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Error binding JSON"})
		return
	}
	status := api.Upload(username, json, doc_id)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(201, gin.H{})
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
	return "", 200
}

