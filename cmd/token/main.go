package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

var tokens [] models.Token
func main() {
	tokens = [] models.Token{}
	api := gin.Default()
	api.POST("/token", createToken)
	api.Run(":8082")
}

func createToken(c *gin.Context) {
	status, username, message := checkBody_Username(c)
	if !status {
		c.JSON(400, gin.H{"error": message})
		return
	}
	token, err := models.CreateToken(username["username"])
	if err != constants.OK {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Error creating token"})
		return
	}
	log.Printf("Token: %s", token.TOKEN)
	tokens = append(tokens, token)
	c.JSON(200, gin.H{"token": token.TOKEN})
}

func checkBody_Username(c *gin.Context) (bool, map[string]string, string) {
	var username map[string]string
	if c.BindJSON(&username) != nil {
		return false, username, "Error binding JSON"
	}
	if username["username"] == "" {
		return false, username, "Username is required"
	}
	return true, username, ""
}