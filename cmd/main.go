package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/models"
	"github.com/luideoz/API_Go/config"
	"encoding/json"
)

func main() {
	api := gin.Default()
	config.Load("config/config.toml")
	api.GET("/version", getVersion)
	api.POST("/signup", signup)
	api.GET("/login", login)
	api.POST("/:username/:doc_id", upload)
	log.Println("Iniciando HTTPS en puerto 8080.....")
	api.Run(":8080")

}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	log.Println("GET /version")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func signup(c *gin.Context) {
	log.Println("POST /signup")
	valid, user,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	user.USERNAME = api.Encrypt_hash(user.USERNAME)
	user.PASSWORD = api.Encrypt_hash(user.PASSWORD)

	/*preparamos el cuerpo json de una request*/
	body, err := json.Marshal(user)
	if err != nil {
		log.Printf("error en la creacion del json %v",err)
		c.JSON(500, gin.H{"error":"error en la creacion del json"})
		return
	}
	response, err := http.Post("http://localhost:8081/signup", "application/json",bytes.NewBuffer(body))
	if err != nil {
		log.Printf("connection error %v",err)
		c.JSON(500, gin.H{"error":"Connection error"})
		return
	}
	defer response.Body.Close()
	var result map[string]string
	json.NewDecoder(response.Body).Decode(&result)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, gin.H{"error":result["error"]})
		return
	}
	c.JSON(201, gin.H{"token":result["token"]})
}
func login(c *gin.Context) {
	valid, user, msg := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":msg})
		return
	}
	user.USERNAME = api.Encrypt_hash(user.USERNAME)
	user.PASSWORD = api.Encrypt_hash(user.PASSWORD)
	body, err := json.Marshal(user)
	if err != nil {
		log.Printf("error en la creacion del json %v",err)
		c.JSON(500, gin.H{"error":"error en la creacion del json"})
		return
	}
	response, err := http.Post("http://localhost:8081/login", "application/json",bytes.NewBuffer(body))
	if err != nil {
		log.Printf("connection error %v",err)
		c.JSON(500, gin.H{"error":"Connection error"})
		return
	}
	defer response.Body.Close()
	var result map[string]string
	json.NewDecoder(response.Body).Decode(&result)
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, gin.H{"error":result["error"]})
		return
	}
	c.JSON(200, gin.H{"token":result["token"]})
}

func upload(c *gin.Context) {
	username := c.Param("username")
	doc_id := c.Param("doc_id")
	username = api.Encrypt_hash(username)
	valid, json_file, msg := checkBody_file(c)
	if !valid {
		c.JSON(400, gin.H{"error":msg})
		return
	}
	valid, token := checkHeader(c)
	if !valid {
		c.JSON(400, gin.H{"error":"empty token"})
		return
	}
	response, err := http.Get("http://localhost:8082/check/"+username+"/"+token)
	if err != nil {
		log.Printf("connection error %v",err)
		c.JSON(500, gin.H{"error":"Connection error"})
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 204 {
		c.JSON(response.StatusCode, gin.H{"error":"invalid token"})
		return
	}
	body,err := json.Marshal(json_file)
	if err != nil {
		c.JSON(500, gin.H{"error":"error en la creacion del json"})
		return
	}
	response, err = http.Post("http://localhost:8083/"+username+"/"+doc_id, "application/json",bytes.NewBuffer(body))
	if err != nil {
		log.Printf("connection error %v",err)
		c.JSON(500, gin.H{"error":"Connection error"})
		return
	}
	defer response.Body.Close()
	var result map[string]string
	json.NewDecoder(response.Body).Decode(&result)
	if response.StatusCode != 201 {
		c.JSON(response.StatusCode, gin.H{"error":result["error"]})
		return
	}
	c.JSON(201, gin.H{"size":result["size"]})
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

func checkHeader(c *gin.Context) (bool,string) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return false, ""
	}
	return true, token
}

func checkBody_file(c *gin.Context) (bool,models.Json,string) {
	var json models.Json
	if c.Request.Body == nil {
		return false, json, "empty body"
	}
	err := c.BindJSON(&json)
	if err != nil {
		return false, json, "invalid json"
	}
	if json.Doc_content == "" {
		return false, json, "empty doc_content"
	}
	return true, json, ""
}
