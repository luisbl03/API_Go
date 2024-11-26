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
	api.GET("/login", login)
	api.POST("/:username/:doc_id", upload)
	api.GET("/:username/:doc_id", getFile)
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
	status = api.Root(username)
	if status == constants.NOT_FOUND {
		c.JSON(404, gin.H{"error":"user not found"})
		return
	}
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error":"internal error (root)"})
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

func login(c *gin.Context) {
	log.Println("GET /login")
	valid, json,message := checkBody_user(c)
	if !valid {
		c.JSON(400, gin.H{"error":message})
		return
	}
	username := json["username"]
	password := json["password"]
	status := api.Login(username,password)
	if status == constants.NOT_FOUND {
		c.JSON(404, gin.H{"error":"user not found"})
		return
	}
	if status == constants.ERROR {
		c.JSON(500,gin.H{"error":"internal error (login)"})
		return
	}
	token, status := models.CreateToken(username)
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error":"internal error (token)"})
		return
	}
	tokens = append(tokens, token)
	c.JSON(200, gin.H{"token":token.TOKEN})
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
	doc_content := json["doc_content"]
	username := c.Param("username")
	doc_id := c.Param("doc_id")
	log.Println("token: ", token)
	log.Println("doc_content: ", doc_content)
	valid = checkToken(token, username)
	if !valid {
		c.JSON(401, gin.H{"error":"unauthorized (invalid token)"})
		return
	}
	status := api.Upload(username, doc_content, doc_id)
	if status == constants.ERROR {
		c.JSON(500, gin.H{"error":"internal error (upload)"})
		return
	}
	if status == constants.EXISTS {
		c.JSON(409, gin.H{"error":"document exists"})
		return
	}
	c.JSON(201, gin.H{"size":status})
}

func getFile(c *gin.Context) {

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

func checkBody_file(c *gin.Context) (bool,map[string]string,string) {
	var json map[string]string
	if c.Request.Body == nil {
		return false, json, "empty body"
	}
	err := c.BindJSON(&json)
	if err != nil {
		return false, json, "invalid json"
	}
	if json["doc_content"] == "" {
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

func checkToken(token string, username string) bool {
	for _, t := range tokens {
		if t.User == username {
			if t.TOKEN == token {
				log.Println("token found: ", t.TOKEN)
				return true
			}
		}
	}
	return false
}