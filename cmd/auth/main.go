package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/models"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/config"
)

var tokens [] models.Token
func main() {
	config.Load("config/config.toml")
	tokens = [] models.Token{}
	api := gin.Default()
	api.POST("/signup", signup)
	api.POST("/login", login)
	api.GET("/:username/:token", checkToken)


	api.Run("localhost:8081")
}

func signup(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error":"invalid json"})
		return
	}
	status := api.Register(user)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	status = api.Root(user.USERNAME)
	if status != constants.OK {
		c.JSON(500, gin.H{"error":"internal error (root)"})
		return
	}
	token, status := models.CreateToken(user.USERNAME)
	if status != constants.OK {
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(201, gin.H{"token":token.TOKEN})
}

func login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error":"invalid json"})
		return
	}
	status := api.Login(user)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	deleteToken(user.USERNAME) //reiniciamos sesion y asi evitamos conflictos con los tokens
	token, status := models.CreateToken(user.USERNAME)
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(200, gin.H{"token":token.TOKEN})
}

func checkToken(c *gin.Context) {
	username := c.Param("username")
	token_code := c.Param("token")
	for _, t := range tokens {
		if t.User == username {
			if t.TOKEN == token_code {
				if models.IsAlive(t) {
					c.JSON(204, gin.H{})
					return
				}
				c.JSON(498, gin.H{"error":"token expired"})
				return
			}
		}
	}
	c.JSON(498, gin.H{"error":"not exists token"})
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

func deleteToken(username string) {
	for i, t := range tokens {
		if t.User == username {
			tokens = append(tokens[:i], tokens[i+1:]...)
			return
		}
	}
}