package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/models"
)

func main() {
	api := gin.Default()
	api.Run(":8083")
}

func upload(c *gin.Context) {
	valid, json, msg := checkBody_file(c)
}

