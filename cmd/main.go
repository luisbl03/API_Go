package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"github.com/luideoz/API_Go/config"
	"github.com/gin-gonic/gin"
)

func main() {
	api := gin.Default()
	config.Load("config/config.toml")
	api.GET("/version", getVersion)
	api.POST("/signup", signup)
	api.Run(":8080")
}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	log.Println("GET /version")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func signup(c *gin.Context) {
	log.Println("POST /signup")
	var request map[string]string
	c.BindJSON(&request)
	jsonData, _ := json.Marshal(request)
	response, err := http.Post("http://localhost:8081/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (request error)"})
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (cant read response)"})
		return
	}
	var status map[string]string
	err = json.Unmarshal(body, &status)
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (cant read status)"})
		return
	}
	log.Println("response: "+string(body))
	code := response.StatusCode
	if code != 201 {
		var msg map[string]string
		json.Unmarshal(body, &msg)
		c.JSON(code, gin.H{"error":msg["error"]})
		return
	}
	var token map[string]string
	err = json.Unmarshal(body, &token)
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (cant read token)"})
		return
	}
	c.JSON(201, gin.H{"token ":token["token"]})

}