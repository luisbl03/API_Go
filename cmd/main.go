package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luideiz/API_Go/api"
	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	//"fmt"
)
var tokens [] models.Token
func main() {
	tokens = [] models.Token{}
	api := gin.Default()
	api.GET("/version", getVersion) 
	api.POST("/signup", register)
	api.GET("/login", login)
	api.POST("/:username/:doc_id", upload)
	api.Run(":8080")

}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func register(c* gin.Context) {
	if !checkBody_User(c) {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	var request map[string]string
	c.BindJSON(&request)
	username := request["username"]
	password := request["password"]
	status := api.Register(username, password)
	if status == constants.EXISTS {
		c.JSON(400, gin.H{"error": "User already exists"})
		return
	}
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	token := createToken(username)
	if token.TOKEN == "" {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(200, gin.H{"token": token.TOKEN})

}

func login(c* gin.Context) {
	if !checkBody_User(c) {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	var request map[string]string
	c.BindJSON(&request)
	username := request["username"]
	password := request["password"]
	status := api.Login(username, password)
	if status == constants.NOT_FOUND {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if status == constants.INVALID_PASSWORD {
		c.JSON(400, gin.H{"error": "Invalid Password"})
		return
	}
	token := createToken(username)
	if token.TOKEN == "" {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(200, gin.H{"token": token.TOKEN})
}

func upload(c* gin.Context) {
	if !checkBody_json(c) {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	if !checkHeader(c) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	token_id := c.Request.Header.Get("Authorization")
	token := getToken(token_id)
	username := c.Param("username")
	if !checkToken(token, username) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
	}
	if !checkExpiration(token) {
		c.JSON(401, gin.H{"error": "Expired Token"})
	}
	doc_id := c.Param("doc_id")
	var data map[string]string
	c.BindJSON(&data)
	doc_content := data["doc_content"]
	status := api.Upload(username, doc_content, doc_id)
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	if status == constants.EXISTS {
		c.JSON(400, gin.H{"error": "Document already exists"})
		return
	}
	c.JSON(200,gin.H{"size":status})
}


func checkBody_User(c *gin.Context) bool {
	if c.Request.Body == nil {
		return false
	}
	//miramos si esta el username y el password
	var json map[string]string
	err := c.BindJSON(&json)
	if err != nil {
		return false
	}
	if json["username"] == "" || json["password"] == "" {
		return false
	}
	return true

}

func checkBody_json(c *gin.Context) bool {
	if c.Request.Body == nil {
		return false
	}
	var body map[string]string
	err := c.BindJSON(&body)
	if err != nil {
		return false
	}
	if body["doc_content"] == "" {
		return false
	}
	return true
}

func checkHeader(c *gin.Context) bool {
	if c.Request.Header.Get("Authorization") == "" {
		return false
	}
	return true
}

func getToken(token string) models.Token {
	for _, t := range tokens {
		if t.TOKEN == token {
			return t
		}
	}
	return models.Token{}
}

func checkToken(token models.Token, user string) bool {
	if token.User != user {
		return false
	}
	if !models.IsAlive(token) {
		return false
	}
	return true
}

func checkExpiration(token models.Token) bool {
	if !models.IsAlive(token) {
		//lo borramos de la lista de tokens
		for i, t := range tokens {
			if t.TOKEN == token.TOKEN {
				tokens = append(tokens[:i], tokens[i+1:]...)
				break
			}
		}
		return false
	}
	return true
}

func createToken(username string) models.Token {
	token, status := models.CreateToken(username)
	if status == constants.ERROR {
		return models.Token{}
	}
	return token
}