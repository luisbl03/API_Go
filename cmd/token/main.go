package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

var tokens []models.Token
func main() {
	tokens = []models.Token{}
	api := gin.Default()
	api.POST("/token", createToken)
	api.Run(":8082")
}


//@Summary Create a token
//@Description Create a token for the user
//@Tags token
//@Accept json
//@Produce json
func createToken(c *gin.Context) {
	//obtenemos el username del cuerpo de la request
	var username map[string]string
	err := c.BindJSON(&username)
	if err != nil {
		log.Println("invalid json")
		c.JSON(400, gin.H{"error":"invalid json"})
		return
	}
	token, status := models.CreateToken(username["username"])
	if status != constants.OK {
		log.Println("internal error (token)")
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(201, gin.H{"token":token.TOKEN})
}