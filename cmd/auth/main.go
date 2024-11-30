package main

import (
	"encoding/json"
	"bytes"
	"io"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

func main() {
	api := gin.Default()
	config.Load("config/config.toml")
	api.POST("/auth/signup", signup)
	api.Run(":8081")

}
func signup(c *gin.Context) {
	log.Println("POST /auth/signup")
	valid, user,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	log.Println("user: "+user.USERNAME)
	log.Println("pass: "+user.PASSWORD)
	user.USERNAME = api.Encrypt_hash(user.USERNAME)
	user.PASSWORD = api.Encrypt_hash(user.PASSWORD)
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
	jsonData := []byte(`{"username":"`+user.USERNAME+`"}`)
	response,err := http.Post("http://localhost:8082/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (500 token)"})
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (cant read token)"})
		return
	}
	/*extraemos el token de la respuesta*/
	var token map[string]string
	err = json.Unmarshal(body, &token)
	if err != nil {
		c.JSON(500, gin.H{"error":"internal error (cant parse token)"})
		return
	}
	log.Println("token: "+token["token"])
	c.JSON(201, gin.H{"token":token["token"]})
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