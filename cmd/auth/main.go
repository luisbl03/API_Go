package main

import (
	"bytes"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

func main() {
	api := gin.Default()
	config.Load("config/config.toml")
	api.POST("/signup", signup)
	api.Run(":8081")
}

func signup(c *gin.Context) {
	log.Println("POST /signup")
	valid, user,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	status := api.Register(user)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	api.Root(user.USERNAME)
	body := map[string]string{"username":user.USERNAME}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error":"Error creating token (json)"})
		return
	}
	response, err := http.Post("http://localhost:8082/token", "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error":"Error creating token (post)"})
		return
	}
	var token map[string]string
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error":"Error creating token (decode)"})
		return
	}
	log.Println("Token: ", token["token"])
	c.JSON(201, gin.H{"token": token["token"]})
}

func checkBody_user(c *gin.Context) (bool,models.User,string) {
	//miramos si el cuerpo del mensaje esta vacio
	var user models.User
	if c.Request.Body == nil {
		return false,user,"empty body"
	}
	err := c.BindJSON(&user)
	if err != nil {
		return false,user,"invalid json"
	}
	//miramos si el campo username esta vacio
	if user.USERNAME == "" {
		return false,user,"empty username"
	}
	//miramos si el campo password esta vacio
	if user.PASSWORD == "" {
		return false,user,"empty password"
	}
	return true,user,""
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