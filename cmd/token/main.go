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
	api.GET("/check/:user/:token", check)
	api.Run(":8082")
}

func createToken(c *gin.Context) {
	status, username, message := checkBody_Username(c)
	if !status {
		c.JSON(400, gin.H{"error": message})
		return
	}
	deleteToken(username["username"])
	token, err := models.CreateToken(username["username"])
	if err != constants.OK {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Error creating token"})
		return
	}
	log.Printf("Token: %s", token.TOKEN)
	tokens = append(tokens, token)
	c.JSON(204, gin.H{"token": token.TOKEN})
}

func check(c *gin.Context) {
	token := c.Param("token")
	username := c.Param("user")
	msg := checkToken(token, username)
	if msg != "" {
		c.JSON(401, gin.H{"error": msg})
		return
	}
	c.JSON(200, gin.H{"status": "OK"})
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
func checkToken(token string, username string) string {
	for _, t := range tokens {
		if t.User == username {
			if t.TOKEN == token {
				if models.IsAlive(t) {
					return ""
				}
				return "expired token"
			}
		}
	}
	return "not found"
}
func deleteToken(username string) {
	for i, t := range tokens {
		if t.User == username {
			tokens = append(tokens[:i], tokens[i+1:]...)
			return
		}
	}
}