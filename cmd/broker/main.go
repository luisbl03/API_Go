package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luideoz/API_Go/api"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/models"
    "io"
)

func main() {
    config.Load("config/config.toml")
    api := gin.Default()
    api.GET("/version", version)
    api.POST("/signup", signup)
    api.GET("/login", login)
    api.POST("/:username/:doc_id", upload)

    api.Run(":8080")
}

func version(c *gin.Context) {
    c.JSON(200, gin.H{
        "version": "1.0.0",
    })
}


func signup(c *gin.Context) {
    valid, user, message := checkBody_user(c)
    if !valid {
        c.JSON(400, gin.H{"error": message})
        return
    }

    user.USERNAME = api.Encrypt_hash(user.USERNAME)
    user.PASSWORD = api.Encrypt_hash(user.PASSWORD)

    jsonData, err := json.Marshal(user)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json marshal"})
        return
    }
    response, err := http.Post("http://localhost:8081/signup", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to auth service"})
        return
    }
    defer response.Body.Close()
    var data map[string]string
    err = json.NewDecoder(response.Body).Decode(&data)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json decode"})
        return
    }
    c.JSON(response.StatusCode, data)
}

func login(c *gin.Context) {
    valid, user, message := checkBody_user(c)
    if !valid {
        c.JSON(400, gin.H{"error": message})
        return
    }

    user.USERNAME = api.Encrypt_hash(user.USERNAME)
    user.PASSWORD = api.Encrypt_hash(user.PASSWORD)

    jsonData, err := json.Marshal(user)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json marshal"})
        return
    }
    response, err := http.Post("http://localhost:8081/login", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to auth service"})
        return
    }
    defer response.Body.Close()
    var data map[string]string
    err = json.NewDecoder(response.Body).Decode(&data)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json decode"})
        return
    }
    c.JSON(response.StatusCode, data)
}

func upload(c *gin.Context) {
    valid, json_data, msg := checkBody_file(c)
    if !valid {
        c.JSON(400, gin.H{"error": msg})
        return
    }
    valid, token := checkHeader(c)
    if !valid {
        c.JSON(400, gin.H{"error": "empty token"})
        return
    }
    username := c.Param("username")
    username = api.Encrypt_hash(username)
    doc_id := c.Param("doc_id")
    response, err := http.Get("http://localhost:8081/"+username+"/"+token)
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to auth service"})
        return 
    }
    defer response.Body.Close()
    if response.StatusCode != 204 {
        body, _ := io.ReadAll(response.Body)
        c.JSON(response.StatusCode, gin.H{"error": string(body)})
    }
    jsonData, err := json.Marshal(json_data)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal error, json marshal"})
        return
    }
    response, err = http.Post("http://localhost:8082/"+username+"/"+doc_id, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        c.JSON(500, gin.H{"error": "could not connect to file service"})
        return
    }
    defer response.Body.Close()
    var data map[string]int
    err = json.NewDecoder(response.Body).Decode(&data)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(response.StatusCode, data)
   
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

func checkHeader(c *gin.Context) (bool,string) {
    token := c.Request.Header.Get("Authorization")
    if token == "" {
        return false, ""
    }
    return true, token
}