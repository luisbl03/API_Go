package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/luideiz/API_Go/api"
	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/config"
	//"fmt"
)

var tokens [] models.Token
func main() {
	tokens = [] models.Token{}
	api := gin.Default()
	config.Load("config/config.toml")
	api.GET("/version", getVersion) 
	api.POST("/signup", register)
	api.GET("/login", login)
	api.POST("/:username/:doc_id", upload)
	api.GET("/:username/:doc_id", getFile)
	api.PUT("/:username/:doc_id", update)
	api.DELETE("/:username/:doc_id", delete)
	api.Run(":8080")

}

func getVersion(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
	log.Println("GET /version")
	c.JSON(200, gin.H{"version":"1.0.0"})
}

func register(c *gin.Context) {
	log.Println("POST /signup")
	valid, user,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}

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
	token, status := models.CreateToken(user.USERNAME)
	if status != constants.OK {
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(code, gin.H{"token":token.TOKEN})
}

func login(c *gin.Context) {
	log.Println("GET /login")
	valid, user,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	user.USERNAME = api.Encrypt_hash(user.USERNAME)
	user.PASSWORD = api.Encrypt_hash(user.PASSWORD)
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
	c.JSON(code, gin.H{"token":token.TOKEN})
}

func upload(c *gin.Context) {
	log.Println("POST /:username/:doc_id")
	valid, json,message := checkBody_file(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	valid, token := checkHeader(c)
	if !valid {
		c.JSON(401, gin.H{"error":"unauthorized (no header)"})
		return
	}
	username := c.Param("username")
	username = api.Encrypt_hash(username)
	doc_id := c.Param("doc_id")
	log.Println("token: ", token)
	log.Println("doc_content: ", json.Doc_content)
	msg := checkToken(token, username)
	if msg != "" {
		c.JSON(401, gin.H{"error":msg})
		return
	}
	status := api.Upload(username, json, doc_id)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(201, gin.H{"size":status})
}

func getFile(c *gin.Context) {
	log.Println("GET /:username/:doc_id")
	username := c.Param("username")
	username = api.Encrypt_hash(username)
	doc_id := c.Param("doc_id")
	valid, token := checkHeader(c)
	if !valid {
		c.JSON(401, gin.H{"error":"unauthorized (no header)"})
		return
	}
	msg := checkToken(token, username)
	if msg != "" {
		c.JSON(401, gin.H{"error":msg})
		return
	}
	json, status := api.GetFile(doc_id, username)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(200, gin.H{"doc_content":json.Doc_content})
}

func update(c *gin.Context) {
	log.Println("PUT /:username/:doc_id")
	username := c.Param("username")
	username = api.Encrypt_hash(username)
	doc_id := c.Param("doc_id")
	valid, token := checkHeader(c)
	if !valid {
		c.JSON(401, gin.H{"error":"unauthorized (no header)"})
		return
	}
	msg := checkToken(token, username)
	if msg != "" {
		c.JSON(401, gin.H{"error":msg})
		return
	}
	valid, json ,msg:= checkBody_file(c)
	if !valid {
		c.JSON(400, gin.H{"error":msg})
		return
	}
	status := api.Update(doc_id, username, json)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(200, gin.H{"size":status})

}

func delete(c *gin.Context) {
	log.Println("DELETE /:username/:doc_id")
	username := c.Param("username")
	username = api.Encrypt_hash(username)
	doc_id := c.Param("doc_id")
	valid, token := checkHeader(c)
	if !valid {
		c.JSON(401, gin.H{"error":"unauthorized (no header)"})
		return
	}
	msg := checkToken(token, username)
	if msg != "" {
		c.JSON(401, gin.H{"error":msg})
		return
	}
	status := api.Delete(doc_id, username)
	msg, code := Status(status)
	if msg != "" {
		c.JSON(code, gin.H{"error":msg})
		return
	}
	c.JSON(204, gin.H{})
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