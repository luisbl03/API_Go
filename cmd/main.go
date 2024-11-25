package main

import (
	"log"

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
	//api.GET("/login", login)
	//api.POST("/:username/:doc_id", upload)
	api.Run(":8080")

}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	log.Println("GET /version")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func register(c *gin.Context) {
	log.Println("POST /signup")
	valid, json,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	username := json["username"]
	password := json["password"]
	status := api.Register(username, password)
	if status == constants.EXISTS {
		c.JSON(409, gin.H{"error":"user exists"})
		return
	}
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error":"internal error (register)"})
		return
	}
	token, status := models.CreateToken(username)
	if status != constants.OK {
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(201, gin.H{"token":token.TOKEN})
}

func checkBody_user(c *gin.Context) (bool,map[string]string,string) {
	//miramos si el cuerpo del mensaje esta vacio
	var json map[string]string
	if c.Request.Body == nil {
		return false,json,"empty body"
	}
	err := c.BindJSON(&json)
	if err != nil {
		return false,json,"invalid json"
	}
	//miramos si el campo username esta vacio
	if json["username"] == "" {
		return false,json,"empty username"
	}
	//miramos si el campo password esta vacio
	if json["password"] == "" {
		return false,json,"empty password"
	}
	return true,json,""
}